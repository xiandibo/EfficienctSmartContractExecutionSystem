package followerserver

import (
	"bytes"
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/vadiminshakov/committer/addrtoip"
	"github.com/vadiminshakov/committer/clientpeer"
	"github.com/vadiminshakov/committer/txpool"
	//"github.com/ethereum/go-ethereum/common"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/openzipkin/zipkin-go"
	zipkingrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	log "github.com/sirupsen/logrus"
	"github.com/vadiminshakov/committer/config"
	"github.com/vadiminshakov/committer/coordinator"
	//"github.com/vadiminshakov/committer/db"
	pb "github.com/vadiminshakov/committer/follower"
	"github.com/vadiminshakov/committer/scexecution"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"sync"
	"time"
)

type Option func(server *FollowerServer) error

const (
	TWO_PHASE   = "two-phase"
	THREE_PHASE = "three-phase"
)

// Server holds server instance, node config and connections to followers (if it's a coordinator node)
// server可能是协调者节点，也可能是follower节点
type FollowerServer struct {
	pb.UnimplementedFollowerServer
	Addr                 string
	Followers            []*clientpeer.FollowerClient
	Config               *config.Config
	GRPCServer           *grpc.Server
	//DB                   db.Database
	DBPath               string
	//这两个hook什么意思
	//ProposeHook          func(req *pb.ProposeRequest) bool
	//CommitHook           func(req *pb.CommitRequest) bool
	//NodeCache            *cache.Cache
	Height               uint64
	cancelCommitOnHeight map[uint64]bool
	mu                   sync.RWMutex
	Tracer               *zipkin.Tracer

	ConfigOfChain		 *scexecution.Config //用于初始化EVM、DB，包括合约状态
	mutest				 sync.Mutex
	TransactionMap		 map[uint64]*Tx

	Txpool				 *txpool.Txpool //交易池

	exitCh	             chan struct{}
	ExitFlag			 bool
}

type Tx struct {
	tx *pb.Transaction
	state bool
}

// InternalCallSmartContract
// 合约内部调用合约的接口
func (s *FollowerServer) InternalCallSmartContract(ctx context.Context, in *pb.InputForSC) (*pb.ResultFromEVM, error) {
	// log.Info("处理其他节点发送来的智能合约调用请求 == FollowerServer.InternalCallSmartContract====", in.Txid)
	//log.Info("本机地址:", s.Addr)

	//s.mutest.Lock() // 防止并发读写statedb
	
	var res1 []byte
	var res2 []byte
	var w sync.WaitGroup
	w.Add(1)
	executeSC(&res1, &res2, s.Addr, in, &w)
	

	statecopy := s.ConfigOfChain.State.CopyStateDBForSubFollower(common.BytesToAddress(in.Address))
	//statecopy = s.ConfigOfChain.State.Copy()
	//log.Info("(Before) length of stateObjects in stateCopy:", statecopy.GetLengthOfStateObjects())

	res, _, err := scexecution.ExecuteTestContract(
		s.ConfigOfChain,
		in.Input,
		common.BytesToAddress(in.Address),
		in.Txid,
		statecopy,
	)
	//log.Info("(After) length of stateObjects in stateCopy:", statecopy.GetLengthOfStateObjects())

	//s.mutest.Unlock()

	w.Wait() // 等待所有Voter执行完成
	//fmt.Println("addr", common.BytesToAddress(in.Address), in.Input)
	//fmt.Println("subCall", res, err)

	if err != nil {
		return nil, err
	}

	//fmt.Println(res)
	//fmt.Println(res1)
	//fmt.Println(res2)

	// 判断执行结果是否互相相等
	// 实际测试合约中set函数没有返回值，也可以匹配执行是否成功
	if bytes.Compare(res, res1) == 0 && bytes.Compare(res, res2) == 0 {
		//fmt.Println("达成共识")
	}

	// 2021/8/29 优化不需要precommit阶段的response
	//log.Info(res)
	//go response(in.Txid) //给协调者发送PreCommit阶段的响应

	return &pb.ResultFromEVM{
		Result: res,
		Issuccess: true,
	}, nil
}

// 2个Voter
// 以协程的方式远程调用, 找到leader的Voter并远程调用
func executeSC(res1, res2 *[]byte, leader string, in *pb.InputForSC, w *sync.WaitGroup)  {
	defer w.Done()

	go func() {
		conn1 := addrtoip.NodeAddrConnMap[addrtoip.LeaderToVoter[leader][0]]
		scpeer1 := &clientpeer.FollowerClient{Connection: pb.NewFollowerClient(conn1)}
		r1, _ := scpeer1.InternalCallSmartContractToVoter(context.Background(), in)
		*res1 = r1.Result
	}()

	go func() {
		conn2 := addrtoip.NodeAddrConnMap[addrtoip.LeaderToVoter[leader][1]]
		scpeer2 := &clientpeer.FollowerClient{Connection: pb.NewFollowerClient(conn2)}
		r2, _ := scpeer2.InternalCallSmartContractToVoter(context.Background(), in)
		*res2 = r2.Result
	}()

}

// 内部共识调用，leader调用voter
func (s *FollowerServer) InternalCallSmartContractToVoter(ctx context.Context, in *pb.InputForSC) (*pb.EVMResultFromVoter, error) {
	//log.Info("共识请求 == InternalCallSmartContractToVoter====", in.Txid)
	//log.Info("本机地址:", s.Addr)

	//s.mutest.Lock() // 防止并发读写statedb

	statecopy := s.ConfigOfChain.State.CopyStateDBForSubFollower(common.BytesToAddress(in.Address))
	//statecopy = s.ConfigOfChain.State.Copy()
	//log.Info("(Before) length of stateObjects in stateCopy:", statecopy.GetLengthOfStateObjects())

	res, _, err := scexecution.ExecuteTestContract(
		s.ConfigOfChain,
		in.Input,
		common.BytesToAddress(in.Address),
		in.Txid,
		statecopy,
	)
	//log.Info("(After) length of stateObjects in stateCopy:", statecopy.GetLengthOfStateObjects())

	//s.mutest.Unlock()



	if err != nil {
		return nil, err
	}

	// 默认达成共识
	s.ConfigOfChain.State.FinaliseFor2PC(common.BytesToAddress(in.Address))
	// 2021/8/29 优化不需要precommit阶段的response
	//log.Info(res)
	//go response(in.Txid) //给协调者发送PreCommit阶段的响应

	// 暂时没有签名
	// TODO：如果要签名，在这里加签名

	return &pb.EVMResultFromVoter{
		Result: res,
		Issuccess: true,
	}, nil
}

// response
// 使得调用返回和发起相应异步执行
func response(txid uint64) {
	coorcli, _ := clientpeer.NewCoordinatorClientPeer(addrtoip.CoordinatorAddr, nil)
	coorcli.PreCommitPhaseResponse(context.Background(), &coordinator.ACK{
		Issuccess:           true,
		Nodeaddr:            nil,
		Txid:                txid,
		PositionOfCallChain: 0,
	})
	//log.Info("InternalCallSmartContract 执行成功，返回PreCommitPhaseResponse， TxID：", txid)
}

// PreCommitRequest
// 来自coordinator的调用，执行主合约，返回执行结果
// 2PC流程的第一步
// TODO：执行结果的形式？应当返回执行结果的状态，ModifyList的形式记录修改了哪些数据。
func (s *FollowerServer) PreCommitRequest(ctx context.Context, in *pb.Transaction) (*emptypb.Empty, error) {
	//事务执行主合约入口
	//log.Info("处理协调者发送来的智能合约调用请求 == FollowerServer.PreCommitRequest====", in.Id)

	accessList := make([]common.Address, 0)
	//解析传过来的访问列表
	for i := 0; i < len(in.AccessAddressList) / 20; i++ {
		accessList = append(accessList, common.BytesToAddress(in.AccessAddressList[i*20: ((i+1)*20)]))
	}
	//log.Info(accessList[0])

	// TODO:不使用粗粒度的锁申请的必要前提是：协调者保证发送的交易不会冲突访问同样的合约
	//s.mutest.Lock()	//粗粒度的加锁，防止并发读写stateDB导致panic

	// 传入本次事务涉及到的合约访问列表，生成一个包含这些合约state object的statedb副本
	statecopy := s.ConfigOfChain.State.CopyStateDBForMainFollower(accessList)
	//statecopy = s.ConfigOfChain.State.Copy()
	//log.Info("(Before) length of stateObjects in stateCopy:", statecopy.GetLengthOfStateObjects())
	// TODO：后续应当改为细粒度的加锁方案，目前的方案存在死锁的风险！死锁情况：本次调用再调用本节点的合约就会死锁！
	_, _, err := scexecution.ExecuteTestContract(
		s.ConfigOfChain,
		in.Input,
		accessList[0], //主合约地址
		in.Id,
		statecopy,
	)
	//s.mutest.Unlock()
	//log.Info("(After) length of stateObjects in stateCopy:", statecopy.GetLengthOfStateObjects())

	//cryp.Sign()

	if err != nil {
		return nil, err
	}
	//返回结果
	//log.Info(res)

	// 2021/8/29 优化不需要precommit阶段的response
	//到这里就完成了本事务的调用链
	//执行完整个事务，向协调者发送Response
	//coorcli, _ := clientpeer.NewCoordinatorClientPeer(addrtoip.CoordinatorAddr, nil)
	//coorcli.PreCommitPhaseResponse(context.Background(), &coordinator.ACK{
	//	Issuccess:           true,
	//	Nodeaddr:            nil,
	//	Txid:                in.Id,
	//	PositionOfCallChain: 0,
	//})
	//log.Info("Main Contract 执行完毕，返回PreCommitPhaseResponse")

	return &empty.Empty{}, nil
}

// CommitRequest 协调者对所有本事务涉及到的follower发送Commit请求
func (s *FollowerServer) CommitRequest(ctx context.Context, in *pb.Transaction) (*emptypb.Empty, error) {
	//持久化本地状态，并返回确认响应
	//本项目的并发粒度最好是合约级别，不是合约里面的数据项

	//每一个事务启动前的state如何生成？提交阶段如何持久化？

	if len(in.AccessAddressList) % 20 != 0 {
		return &emptypb.Empty{}, nil
	}

	accessList := make([]common.Address, 0)

	//解析传过来的访问列表
	//TODO：这个是重复的步骤，应当直接传送合约地址，不需要每次都解析一遍
	for i := 0; i < len(in.AccessAddressList) / 20; i++ {
		accessList = append(accessList, common.BytesToAddress(in.AccessAddressList[i*20: ((i+1)*20)]))
	}

	scAddr := accessList[in.PositionOfCallChain]

	//s.mutest.Lock()
	s.ConfigOfChain.State.FinaliseFor2PC(scAddr)
	//s.mutest.Unlock()

	//log.Info("CommitRequest====", in.Id)
	//完成持久化，给协调者发送相应，调用协调者响应方法s
	go CommitResToCoor(in)

	return &empty.Empty{}, nil
}

func CommitResToCoor(in *pb.Transaction)  {
	coorcli, _ := clientpeer.NewCoordinatorClientPeer(addrtoip.CoordinatorAddr, nil)
	coorcli.CommitPhaseResponse(context.Background(), &coordinator.ACK{
		Issuccess:           true,
		Nodeaddr:            nil,
		Txid:                in.Id,
		PositionOfCallChain: in.PositionOfCallChain,
	})
	//log.Info("CommitResToCoor CommitRequest in follower to Coordinator, TxID:", in.Id)
}








// WithBadgerDB adds BadgerDB manager to the Server instance
//func WithBadgerDB(path string) func(*FollowerServer) error {
//	return func(server *FollowerServer) error {
//		var err error
//		server.DB, err = db.New(path)
//		return err
//	}
//}

/*func checkServerFields(server *FollowerServer) error {
	if server.DB == nil {
		return errors.New("database is not selected")
	}
	return nil
}*/

// Run starts non-blocking GRPC server
func (s *FollowerServer) Run(opts ...grpc.UnaryServerInterceptor) {
	var err error

	//创建GRPC服务器
	if s.Config.WithTrace {
		s.GRPCServer = grpc.NewServer(grpc.ChainUnaryInterceptor(opts...), grpc.StatsHandler(zipkingrpc.NewServerHandler(s.Tracer)))
	} else {
		s.GRPCServer = grpc.NewServer(grpc.ChainUnaryInterceptor(opts...))
	}
	pb.RegisterFollowerServer(s.GRPCServer, s)

	//监听本节点地址端口
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Infof("listening on tcp://%s", s.Addr)

	go s.GRPCServer.Serve(l)
}

// Stop stops server
func (s *FollowerServer) Stop() {
	log.Info("stopping server")
	s.GRPCServer.GracefulStop()
	//if err := s.DB.Close(); err != nil {
	//	log.Infof("failed to close db, err: %s\n", err)
	//}
	log.Info("server stopped")
}

//func (s *FollowerServer) rollback() {
//	s.NodeCache.Delete(s.Height)
//}

func (s *FollowerServer) SetProgressForCommitPhase(height uint64, docancel bool) {
	s.mu.Lock()
	s.cancelCommitOnHeight[height] = docancel
	s.mu.Unlock()
}

func (s *FollowerServer) HasProgressForCommitPhase(height uint64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cancelCommitOnHeight[height]
}

// NewFollowerServer fabric func for Server
func NewFollowerServer(conf *config.Config, opts ...Option) (*FollowerServer, error) {
	//writer3, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE, 0766)

	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true, // Seems like automatic color detection doesn't work on windows terminals
		FullTimestamp:   true,
		TimestampFormat: time.StampMicro,
	})

	//创建sever地址
	server := &FollowerServer{Addr: conf.Nodeaddr, DBPath: conf.DBPath}
	var err error
	//绑定propose 和 commit
	for _, option := range opts {
		err = option(server)
		if err != nil {
			return nil, err
		}
	}

	//if conf.WithTrace {
	//	// get Zipkin tracer
	//	server.Tracer, err = trace.Tracer(fmt.Sprintf("%s:%s", conf.Role, conf.Nodeaddr), conf.Nodeaddr)
	//	if err != nil {
	//		return nil, err
	//	}
	//}

	//如果输入了follower，初始化FollowerClient
	for _, node := range conf.Followers {
		//cli, err := peer.New(node, server.Tracer)
		cli, err := clientpeer.NewFollowerClientPeer(node, server.Tracer)
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

	//server.DB, err = db.New(conf.DBPath)
	//if err != nil {
	//	return nil, err
	//}
	//server.NodeCache = cache.New()
	//server.cancelCommitOnHeight = map[uint64]bool{}

	if server.Config.CommitType == TWO_PHASE {
		log.Info("two-phase-commit mode enabled")
	} else {
		log.Info("three-phase-commit mode enabled")
	}
	//err = checkServerFields(server)

	//server.TransactionMap = make(map[uint64]*Tx)

	//交易池
	server.Txpool = txpool.NewTxpool()
	server.exitCh = make(chan struct{})

	return server, err
}
