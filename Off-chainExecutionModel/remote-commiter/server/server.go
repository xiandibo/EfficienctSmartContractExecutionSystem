package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/openzipkin/zipkin-go"
	zipkingrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	log "github.com/sirupsen/logrus"
	"github.com/vadiminshakov/committer/cache"
	"github.com/vadiminshakov/committer/config"
	"github.com/vadiminshakov/committer/db"
	"github.com/vadiminshakov/committer/peer"
	pb "github.com/vadiminshakov/committer/proto"
	"github.com/vadiminshakov/committer/scexecution"
	"github.com/vadiminshakov/committer/trace"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	status "google.golang.org/grpc/status"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Option func(server *Server) error

const (
	TWO_PHASE   = "two-phase"
	THREE_PHASE = "three-phase"
)

// Server holds server instance, node config and connections to followers (if it's a coordinator node)
// server可能是协调者节点，也可能是follower节点
type Server struct {
	pb.UnimplementedCommitServer
	Addr                 string
	Followers            []*peer.CommitClient
	Config               *config.Config
	GRPCServer           *grpc.Server
	DB                   db.Database
	DBPath               string
	ProposeHook          func(req *pb.ProposeRequest) bool
	CommitHook           func(req *pb.CommitRequest) bool
	NodeCache            *cache.Cache
	Height               uint64
	cancelCommitOnHeight map[uint64]bool
	mu                   sync.RWMutex
	Tracer               *zipkin.Tracer

	ConfigOfChain		 *scexecution.Config //用于初始化EVM、DB，包括合约状态
}

//Follower被coordinator调用的接口
func (s *Server) Propose(ctx context.Context, req *pb.ProposeRequest) (*pb.Response, error) {
	var span zipkin.Span
	if s.Tracer != nil {
		span, ctx = s.Tracer.StartSpanFromContext(ctx, "ProposeHandle")
		defer span.Finish()
	}
	s.SetProgressForCommitPhase(req.Index, false)
	return s.ProposeHandler(ctx, req, s.ProposeHook)
}

func (s *Server) Precommit(ctx context.Context, req *pb.PrecommitRequest) (*pb.Response, error) {
	var span zipkin.Span
	if s.Tracer != nil {
		span, _ = s.Tracer.StartSpanFromContext(ctx, "PrecommitHandle")
		defer span.Finish()
	}
	if s.Config.CommitType == THREE_PHASE {
		ctx, _ = context.WithTimeout(context.Background(), time.Duration(s.Config.Timeout)*time.Millisecond)
		go func(ctx context.Context) {
		ForLoop:
			for {
				select {
				case <-ctx.Done():
					md := metadata.Pairs("mode", "autocommit")
					ctx := metadata.NewOutgoingContext(context.Background(), md)
					if !s.HasProgressForCommitPhase(req.Index) {
						s.Commit(ctx, &pb.CommitRequest{Index: s.Height})
						log.Warn("commit without coordinator after timeout")
					}
					break ForLoop
				}
			}
		}(ctx)
	}
	return s.PrecommitHandler(ctx, req)
}

func (s *Server) Commit(ctx context.Context, req *pb.CommitRequest) (resp *pb.Response, err error) {
	var span zipkin.Span
	if s.Tracer != nil {
		span, ctx = s.Tracer.StartSpanFromContext(ctx, "CommitHandle")
		defer span.Finish()
	}

	if s.Config.CommitType == THREE_PHASE {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Internal, "no metadata")
		}

		meta := md["mode"]

		if len(meta) == 0 {
			s.SetProgressForCommitPhase(s.Height, true) // set flag for cancelling 'commit without coordinator' action, coz coordinator responded actually
			resp, err = s.CommitHandler(ctx, req, s.CommitHook, s.DB)
			if err != nil {
				return nil, err
			}
			if resp.Type == pb.Type_ACK {
				atomic.AddUint64(&s.Height, 1)
			}
			return
		}

		if req.IsRollback {
			s.rollback()
		}
	} else {

		resp, err = s.CommitHandler(ctx, req, s.CommitHook, s.DB)
		if err != nil {
			return
		}

		if resp.Type == pb.Type_ACK {
			atomic.AddUint64(&s.Height, 1)
		}
		return
	}
	return
}

func (s *Server) Get(ctx context.Context, req *pb.Msg) (*pb.Value, error) {
	var span zipkin.Span
	if s.Tracer != nil {
		span, ctx = s.Tracer.StartSpanFromContext(ctx, "GetHandle")
		defer span.Finish()
	}

	value, err := s.DB.Get(req.Key)
	if err != nil {
		return nil, err
	}
	return &pb.Value{Value: value}, nil
}

//被客户端远程调用，写入一个key
//可以视作一次事务执行过程
//一次2PC流程
func (s *Server) Put(ctx context.Context, req *pb.Entry) (*pb.Response, error) {
	//debug.PrintStack()

	var (
		response *pb.Response
		err      error
		span     zipkin.Span
	)

	if s.Tracer != nil {
		span, ctx = s.Tracer.StartSpanFromContext(ctx, "PutHandle")
		defer span.Finish()
	}

	//我们只考虑2PC
	var ctype pb.CommitType
	if s.Config.CommitType == THREE_PHASE {
		ctype = pb.CommitType_THREE_PHASE_COMMIT
	} else {
		ctype = pb.CommitType_TWO_PHASE_COMMIT
	}

	// propose
	log.Infof("propose key %s", req.Key)
	for _, follower := range s.Followers {
		if s.Config.CommitType == THREE_PHASE {
			ctx, _ = context.WithTimeout(ctx, time.Duration(s.Config.Timeout)*time.Millisecond)
		}
		if s.Tracer != nil {
			span, ctx = s.Tracer.StartSpanFromContext(ctx, "Propose")
		}
		var (
			response *pb.Response
			err      error
		)
		//持续给一个follower发送消息，直到返回消息，缺乏并发性？是不是应该同时给多个follower发消息
		for response == nil || response != nil && response.Type == pb.Type_NACK {
			//给follower发送propose消息
			response, err = follower.Propose(ctx, &pb.ProposeRequest{Key: req.Key,
				Value:      req.Value,
				CommitType: ctype,
				Index:      s.Height})
			if s.Tracer != nil && span != nil {
				span.Finish()
			}
			if response != nil && response.Index > s.Height {
				log.Warnf("сoordinator has stale height [%d], update to [%d] and try to send again", s.Height, response.Index)
				s.Height = response.Index
			}
		}
		if err != nil {
			log.Errorf(err.Error())
			return &pb.Response{Type: pb.Type_NACK}, nil
		}
		//本地写入事务
		s.NodeCache.Set(s.Height, req.Key, req.Value)
		if response.Type != pb.Type_ACK {
			return nil, status.Error(codes.Internal, "follower not acknowledged msg")
		}
	}

	// precommit phase only for three-phase mode
	//这里2PC还是会执行
	log.Infof("precommit key %s", req.Key)
	for _, follower := range s.Followers {
		if s.Config.CommitType == THREE_PHASE {
			ctx, _ = context.WithTimeout(ctx, time.Duration(s.Config.Timeout)*time.Millisecond)
		}
		if s.Tracer != nil {
			span, ctx = s.Tracer.StartSpanFromContext(ctx, "Precommit")
		}
		response, err = follower.Precommit(ctx, &pb.PrecommitRequest{Index: s.Height})
		if s.Tracer != nil && span != nil {
			span.Finish()
		}
		if err != nil {
			log.Errorf(err.Error())
			return &pb.Response{Type: pb.Type_NACK}, nil
		}
		if response.Type != pb.Type_ACK {
			return nil, status.Error(codes.Internal, "follower not acknowledged msg")
		}
	}

	// the coordinator got all the answers, so it's time to persist msg and send commit command to followers
	key, value, ok := s.NodeCache.Get(s.Height)
	if !ok {
		return nil, status.Error(codes.Internal, "can't to find msg in the coordinator's cache")
	}
	if err = s.DB.Put(key, value); err != nil {
		return &pb.Response{Type: pb.Type_NACK}, status.Error(codes.Internal, "failed to save msg on coordinator")
	}

	// commit
	log.Infof("commit %s", req.Key)
	for _, follower := range s.Followers {
		if s.Tracer != nil {
			span, ctx = s.Tracer.StartSpanFromContext(ctx, "Commit")
		}
		response, err = follower.Commit(ctx, &pb.CommitRequest{Index: s.Height})
		if s.Tracer != nil && span != nil {
			span.Finish()
		}
		if err != nil {
			log.Errorf(err.Error())
			return &pb.Response{Type: pb.Type_NACK}, nil
		}
		if response.Type != pb.Type_ACK {
			return nil, status.Error(codes.Internal, "follower not acknowledged msg")
		}
	}
	log.Infof("committed key %s", req.Key)
	// increase height for next round
	atomic.AddUint64(&s.Height, 1)

	return &pb.Response{Type: pb.Type_ACK}, nil
}

func (s *Server) NodeInfo(ctx context.Context, req *empty.Empty) (*pb.Info, error) {
	var span zipkin.Span
	if s.Tracer != nil {
		span, ctx = s.Tracer.StartSpanFromContext(ctx, "NodeInfoHandle")
		defer span.Finish()
	}

	return &pb.Info{Height: s.Height}, nil
}




// 处理其他节点发送来的智能合约调用请求
// 此处为从节点的逻辑，完成本地调用以及调用其他节点的合约
// 还需要创建调用协调者节点的方法
// 被客户端或协调者调用，即发起一个事务
func (s *Server) CallSmartContract(ctx context.Context, input *pb.InputForSC) (*pb.ResultFromEVM, error) {
	fmt.Println("处理其他节点发送来的智能合约调用请求 == CallSmartContract")
	log.Info("本机地址:", s.Addr)

	//初始化EVM
	//需要事先写入合约和数据
	//此调用暂不涉及follower


	//var res []byte
	//res = append(res, byte(11))
	res, _, err := scexecution.ExecuteTestContract(
		s.ConfigOfChain,
		input.Input,
		common.BytesToAddress(input.Address),
		uint64(1),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &pb.ResultFromEVM{
		Result: res,
		Issuccess: true,
	}, nil
}

// 合约调用合约入口
// 接收输入，初始化EVM执行并返回结果
func (s *Server) CallSmartContractInFollower(ctx context.Context, input *pb.InputForSC) (*pb.ResultFromEVM, error) {
	fmt.Println("处理其他节点发送来的智能合约调用请求 == CallSmartContractInFollower")

	//测试的设置，合约的代码状态初始化应放在节点启动时刻
	//scexecution.Init(nil)

	//初始化EVM
	//需要事先写入合约和数据
	//此调用不考虑内部再次内部调用的情况，即默认无内部调用
	res, _, err := scexecution.ExecuteTestContract(
		s.ConfigOfChain,
		input.Input,
		common.BytesToAddress(input.Address),
		uint64(1),
		nil,
	)
	if err != nil {
		return nil, err
	}

	//var res []byte
	//res = append(res, byte(11))
	//res := scexecution.ExecuteSmartContract() //测试执行合约
	return &pb.ResultFromEVM{
		Result: res,
		Issuccess: true,
	}, nil
}

// 客户端调用协调者方法SendTx，协调者执行2PC逻辑
// 此处暂时支持单一调用，以后再设计并发处理TX的逻辑
func (s *Server) SendTx(ctx context.Context, input *pb.Transaction) (*pb.Response, error) {

	//首先解析AccessAddressList
	//每20个字节为一个地址
	if len(input.AccessAddressList) % 20 != 0 {
		return &pb.Response{}, nil
	}

	accessList := make([]common.Address, 0)

	//解析传过来的访问列表
	for i := 0; i < len(input.AccessAddressList) / 20; i++ {
		accessList = append(accessList, common.BytesToAddress(input.AccessAddressList[i*20: ((i+1)*20)]))
	}
	log.Info(accessList)




	return &pb.Response{}, nil
}



// NewCommitServer fabric func for Server
func NewCommitServer(conf *config.Config, opts ...Option) (*Server, error) {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true, // Seems like automatic color detection doesn't work on windows terminals
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})

	//创建sever地址
	server := &Server{Addr: conf.Nodeaddr, DBPath: conf.DBPath}
	var err error
	//绑定propose 和 commit
	for _, option := range opts {
		err = option(server)
		if err != nil {
			return nil, err
		}
	}

	if conf.WithTrace {
		// get Zipkin tracer
		server.Tracer, err = trace.Tracer(fmt.Sprintf("%s:%s", conf.Role, conf.Nodeaddr), conf.Nodeaddr)
		if err != nil {
			return nil, err
		}
	}

	//如果输入了follower，初始化commitclient
	for _, node := range conf.Followers {
		cli, err := peer.New(node, server.Tracer)
		if err != nil {
			return nil, err
		}
		server.Followers = append(server.Followers, cli)
	}

	server.Config = conf
	//如果本节点是协调者，初始化设定本届点协调者的地址
	if conf.Role == "coordinator" {
		server.Config.Coordinator = server.Addr
	}

	server.DB, err = db.New(conf.DBPath)
	if err != nil {
		return nil, err
	}
	server.NodeCache = cache.New()
	server.cancelCommitOnHeight = map[uint64]bool{}

	if server.Config.CommitType == TWO_PHASE {
		log.Info("two-phase-commit mode enabled")
	} else {
		log.Info("three-phase-commit mode enabled")
	}
	err = checkServerFields(server)
	return server, err
}

// WithBadgerDB adds BadgerDB manager to the Server instance
func WithBadgerDB(path string) func(*Server) error {
	return func(server *Server) error {
		var err error
		server.DB, err = db.New(path)
		return err
	}
}

func checkServerFields(server *Server) error {
	if server.DB == nil {
		return errors.New("database is not selected")
	}
	return nil
}

// Run starts non-blocking GRPC server
func (s *Server) Run(opts ...grpc.UnaryServerInterceptor) {
	var err error

	//创建GRPC服务器
	if s.Config.WithTrace {
		s.GRPCServer = grpc.NewServer(grpc.ChainUnaryInterceptor(opts...), grpc.StatsHandler(zipkingrpc.NewServerHandler(s.Tracer)))
	} else {
		s.GRPCServer = grpc.NewServer(grpc.ChainUnaryInterceptor(opts...))
	}
	pb.RegisterCommitServer(s.GRPCServer, s)

	//监听本节点地址端口
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Infof("listening on tcp://%s", s.Addr)

	go s.GRPCServer.Serve(l)
}

// Stop stops server
func (s *Server) Stop() {
	log.Info("stopping server")
	s.GRPCServer.GracefulStop()
	if err := s.DB.Close(); err != nil {
		log.Infof("failed to close db, err: %s\n", err)
	}
	log.Info("server stopped")
}

func (s *Server) rollback() {
	s.NodeCache.Delete(s.Height)
}

func (s *Server) SetProgressForCommitPhase(height uint64, docancel bool) {
	s.mu.Lock()
	s.cancelCommitOnHeight[height] = docancel
	s.mu.Unlock()
}

func (s *Server) HasProgressForCommitPhase(height uint64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cancelCommitOnHeight[height]
}
