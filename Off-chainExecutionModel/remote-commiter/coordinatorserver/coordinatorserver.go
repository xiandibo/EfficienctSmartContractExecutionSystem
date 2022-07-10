package coordinatorserver

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/openzipkin/zipkin-go"
	zipkingrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	log "github.com/sirupsen/logrus"
	"github.com/vadiminshakov/committer/addrtoip"
	"github.com/vadiminshakov/committer/cache"
	"github.com/vadiminshakov/committer/clientpeer"
	"github.com/vadiminshakov/committer/config"
	pb "github.com/vadiminshakov/committer/coordinator"
	"github.com/vadiminshakov/committer/follower"
	"github.com/vadiminshakov/committer/scexecution"
	"github.com/vadiminshakov/committer/txpool"
	"github.com/vadiminshakov/committer/utils"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"sync"
	"sync/atomic"
	"time"
)



type Option func(server *CoordinatorServer) error

const (
	TWO_PHASE   = "two-phase"
	THREE_PHASE = "three-phase"
	ACCESS_LIST_LENGTH = 4
)

var CommittedTx uint64 = 0
var TotalTx uint64 = 0
var Start time.Time
var End time.Time

// Server holds server instance, node config and connections to followers (if it's a coordinator node)
// server可能是协调者节点，也可能是follower节点
type CoordinatorServer struct {
	pb.UnimplementedCoordinatorServer
	Addr                 string
	Followers            []*clientpeer.FollowerClient
	Config               *config.Config
	GRPCServer           *grpc.Server
	//DB                   db.Database
	DBPath               string
	//这两个hook什么意思
	//ProposeHook          func(req *pb.ProposeRequest) bool
	//CommitHook           func(req *pb.CommitRequest) bool
	NodeCache            *cache.Cache
	Height               uint64
	cancelCommitOnHeight map[uint64]bool
	mu                   sync.RWMutex
	Tracer               *zipkin.Tracer

	ConfigOfChain		 *scexecution.Config //用于初始化EVM、DB，包括合约状态
	//TransactionMap		 map[uint64]*Tx
	Txpool				 *txpool.Txpool //交易池

	ExitCh	             chan struct{}
	ExitFlag			 bool

	MaxProcessingTx      uint64
	Dispatcher			 *utils.EventDispatcher


}

var NewTxCh = make(chan *txpool.Tx, 10000)

// SendTransaction
// 事务流程的起点
// 1. 客户端构造事务并调用此接口
// 2. 协调者分析AccessList，得到合约调用链
//Client向Coordinator发送交易
// 将tx保存起来
// 为此tx启动一次2pc流程，注意并发
func (s *CoordinatorServer) SendTransaction(ctx context.Context, in *pb.Transaction) (*emptypb.Empty, error) {

	//首先解析AccessAddressList
	//每20个字节为一个地址
	//log.Info("CoordinatorServer.SendTransaction, ", in)
	//不是20个字节的整数倍为不合法访问列表
	//log.Info(len(in.AccessAddressList))
	if len(in.AccessAddressList) % 20 != 0 {
		log.Info("非法合约访问地址")
		return &emptypb.Empty{}, nil
	}
	accessList := make([]common.Address, 0) //保存访问列表
	//解析传过来的访问列表
	for i := 0; i < len(in.AccessAddressList) / 20; i++ {
		accessList = append(accessList, common.BytesToAddress(in.AccessAddressList[i*20: ((i+1)*20)]))
	}
	//log.Info(accessList)

	//新建事务
	tx := &txpool.Tx{
		Tx: in,
		State: false,
		AccessList: accessList,
	}
	// 把事务添加到交易池，TxID重复的事务将会被丢弃（暂定）
	if s.Txpool.AddUnProcessingTx(tx) == false {
		return &empty.Empty{}, nil
	}

	// 2021/8/22 先存交易再启动执行版本

	// 暂时不启动处理事务
	//s.Dispatcher.DispatchEvent("New Transaction", tx) // 分发事件 "New Transaction"

	log.Info("Receive transaction from Client, TransactionID: ", in.Id)
	// 打印
	log.Info("TxPool length:", s.Txpool.TxIn2PCCMap.Count())

	// 启动多个worker线程执行Tx
	// 应该另外写一个接收到命令行输入就启动的函数

	return &empty.Empty{}, nil
}

// CommitPhaseResponse
// 收到follower的响应
// 统计此事务的相应数量，当响应全部收集齐，确认此事务成功执行，其所有修改在对应的节点上都生效
// 此时协调者把此事务的执行结果报告到链上（逐个事务上报还是逐批上报取决于链上链下的具体交互过程），如果把链下事务标记一下与链上事务正常打包进同一个区块中，则最好是逐个发送成功的事务。
func (s *CoordinatorServer) CommitPhaseResponse(ctx context.Context, in *pb.ACK) (*emptypb.Empty, error) {
	//被follower调用
	//TODO: 测试显示在这里堵住了
	if in.Issuccess == true {
		//log.Info("Get Tx", in.Txid)
		tx := s.Txpool.GetUnProcessingTx(in.Txid)
		//log.Info("Get Tx success", in.Txid)
		//原子操作
		//原子操作
		atomic.AddUint64(&tx.ResponseCount, 1)
		if atomic.CompareAndSwapUint64(&tx.ResponseCount, ACCESS_LIST_LENGTH, 0) {
			//log.Info("CommitPhaseResponse 收集齐全部响应, TxID: ", in.Txid)
			// 释放对涉及到的事务的锁
			//log.Info("tx.LockList ======", tx.LockList)
			for _, addr := range tx.LockList {
				s.Txpool.AddrMutexMap[addr].Unlock()
			}

			//发起对该事务的上链流程
			//TODO: 至此单个事务的链下2PC流程已完成，事务上链方式待定
			//可以先存起来，也可以直接作为链下特殊的交易发送到链上
			go CommitTransactionToOnChain(tx)
		}
	}

	return &empty.Empty{}, nil
}

// PreCommitPhaseResponse
// 收集事务PreCommit阶段的响应
func (s *CoordinatorServer) PreCommitPhaseResponse(ctx context.Context, in *pb.ACK) (*emptypb.Empty, error) {
	//其他节点返回的执行结果，应该包含事务标志，节点地址，位于调用链上的位置

	if in.Issuccess == true { //只处理成功执行的事务
		//找到正在处理的事务
		tx := s.Txpool.GetUnProcessingTx(in.Txid)
		//原子操作 +1
		atomic.AddUint64(&tx.ResponseCount, 1)
		if atomic.CompareAndSwapUint64(&tx.ResponseCount, 4, 0) { //因为是每个地址读写各一次，所以是10*2+1=21
			//收集齐响应，启动Commit阶段
			log.Info("PreCommitPhaseResponse All Receive., TxID: ", in.Txid, " Count: ", atomic.LoadUint64(&tx.ResponseCount))

			go CommitPhase(tx) //启动Commit阶段
		}
	} else {
		//TODO:如果有一个事务执行失败，进入事务回滚流程，此处不涉及实验，暂时不考虑

	}

	return &empty.Empty{}, nil
}

// CommitPhase 对事务tx启动提交阶段
// 协调者收集满足够的PrecommitResponse以后，发起CommitPhase
func CommitPhase(tx *txpool.Tx) {
	//这里我们没有每个合约的读写信息，我们还是默认主合约的逻辑是按照顺序对每个合约进行一次读和写
	//所以是对
	//log.Info("对事务启动CommitPhase TxID: ", tx.Tx.Id) //测试提交阶段堵住是否是在这里堵住
	//log.Info("CommitPhase", len(tx.AccessList))
	for _, scaddr := range tx.AccessList {

		nodeaddr := addrtoip.AddressToIP[scaddr]
		cli, _ := clientpeer.NewFollowerClientPeer(nodeaddr, nil)
		tempTx := &follower.Transaction{
			Id:                tx.Tx.Id,
			NumberOfAddress:   tx.Tx.NumberOfAddress,
			AccessAddressList: tx.Tx.AccessAddressList,
			Input:             tx.Tx.Input,
		}
		//注意此处并发操作是否正确？
		//是否因为没有并发才这么慢？
		go cli.CommitRequest(context.Background(), tempTx) // 使用协程是想尽快把所有Commit请求发出去
	}
}


// CommitTransactionToOnChain
// TODO：把成功的事务提交到链上,涉及到链上流程
func CommitTransactionToOnChain(tx *txpool.Tx) {
	// 用于计算全部执行完成的时间
	atomic.AddUint64(&CommittedTx, 1)
	if atomic.CompareAndSwapUint64(&CommittedTx, TotalTx, 0) {
		End = time.Now()
		//fmt.Println("全部执行完成！用时：", End.Sub(Start).Seconds())
		dif := End.Sub(Start).Seconds()
		tps := float64(TotalTx) / dif
		fmt.Printf("全部执行完成！用时: %fs, Tx总数: %d, TPS: %f", dif, TotalTx, tps)
	}
	//log.Info("Commit Transaction to OnChain TxID: ", tx.Tx.Id)
}

// NewCoordinatorServer NewCommitServer fabric func for Server
func NewCoordinatorServer(conf *config.Config, opts ...Option) (*CoordinatorServer, error) {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true, // Seems like automatic color detection doesn't work on windows terminals
		FullTimestamp:   true,
		TimestampFormat: time.StampMicro,
	})

	//创建sever地址
	server := &CoordinatorServer{Addr: conf.Nodeaddr, DBPath: conf.DBPath}
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
	server.ExitCh = make(chan struct{})

	//新建事件监听器
	server.Dispatcher = utils.NewEventDispatcher()

	return server, err
}

// WithBadgerDB adds BadgerDB manager to the Server instance
//func WithBadgerDB(path string) func(*CoordinatorServer) error {
//	return func(server *CoordinatorServer) error {
//		var err error
//		server.DB, err = db.New(path)
//		return err
//	}
//}

/*func checkServerFields(server *CoordinatorServer) error {
	if server.DB == nil {
		return errors.New("database is not selected")
	}
	return nil
}*/

// Run starts non-blocking GRPC server
// 启动协调者服务进程
// 需要关注
// 1. 监听事件调度器，监听 “New Transaction”
// 2. 启动了一个一直等待 “New Transaction” 事件的协程，每收到一个新事务就相应地启动一个2PC流程
func (s *CoordinatorServer) Run(opts ...grpc.UnaryServerInterceptor) {
	var err error

	//创建GRPC服务器
	if s.Config.WithTrace {
		s.GRPCServer = grpc.NewServer(grpc.ChainUnaryInterceptor(opts...), grpc.StatsHandler(zipkingrpc.NewServerHandler(s.Tracer)))
	} else {
		s.GRPCServer = grpc.NewServer(grpc.ChainUnaryInterceptor(opts...))
	}
	pb.RegisterCoordinatorServer(s.GRPCServer, s)

	//监听本节点地址端口
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Infof("listening on tcp://%s", s.Addr)

	go s.GRPCServer.Serve(l)

	//随时检查需要处理的链下交易
	listener := utils.NewEventListener(NewTransactionEventListener)
	s.Dispatcher.AddEventListener("New Transaction", listener)
	go s.ProcessTransactions()
}

//接受一个新事件，表示新的交易到来
func NewTransactionEventListener(event utils.Event) {
	tx := event.Object.(*txpool.Tx) //以txpool.Tx的格式传递事件
	//收到新事务，放入通道
	NewTxCh <- tx
}

// Stop stops server
func (s *CoordinatorServer) Stop() {
	log.Info("stopping server")
	s.GRPCServer.GracefulStop()
	//if err := s.DB.Close(); err != nil {
	//	log.Infof("failed to close db, err: %s\n", err)
	//}
	log.Info("server stopped")
	s.ExitFlag = true
}

func (s *CoordinatorServer) rollback() {
	s.NodeCache.Delete(s.Height)
}

func (s *CoordinatorServer) SetProgressForCommitPhase(height uint64, docancel bool) {
	s.mu.Lock()
	s.cancelCommitOnHeight[height] = docancel
	s.mu.Unlock()
}

func (s *CoordinatorServer) HasProgressForCommitPhase(height uint64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cancelCommitOnHeight[height]
}

// ProcessTransactions
// 监听NewTxCh通道，
//未来应当设定线程上限，服务器能够同时处理多少个事务？是性能瓶颈。
//尝试不用交易池，用通道时间传递事务？
func (s *CoordinatorServer) ProcessTransactions() {



	//等待New Transaction Event，从通道取出新的交易，启动2PC流程
	for {
		select {
		case tx := <- NewTxCh:
			//新事务，发起2PC流程
			// TODO：使用协程与不使用协程两种情况哪个更好？

			// 检查到新事务，进行必要的可串行化检查与操作，满足条件了才发起2PC事务
			// 申请了锁记得释放锁
			// 这样写影响了协调者接收交易的速度
			//if s.Txpool.IsSerializable(tx) {
			//	go s.Arise2PCPreCommitRequest(tx)
			//}
			go s.Arise2PCPreCommitRequest(tx)
		case <- s.ExitCh:
			log.Info("Exit ProcessTransactions()")
			return
		}
	}
}

// Arise2PCPreCommitRequest
// 对一个新的事务发起2PC流程
// 直接对AccessList的第一个合约地址发起调用，我们默认第一个地址是主合约地址
func (s *CoordinatorServer) Arise2PCPreCommitRequest(tx *txpool.Tx) {
	//先申请锁，由于保证不存在死锁，因此肯定能申请到，前提是初始化时包括了所需的合约
	s.Txpool.IsSerializable(tx) // 先申请到相应的锁，保证可串行化

	//第一个地址为主合约地址
	//此处地址默认为共识组的代理节点
	//调用follower的接口
	//为什么原本的2pc实现要持续发送请求直到有返回？是否不保证可靠发送？
	//log.Info(tx.AccessList)
	hostaddr := addrtoip.AddressToIP[tx.AccessList[0]] // hostaddr 例子 "localhos:3001"
	//log.Info(hostaddr)
	cli, err := clientpeer.NewFollowerClientPeer(hostaddr, nil)

	// 远程调用Follower的接口
	_, err = cli.PreCommitRequest(context.Background(), &follower.Transaction{
		Id:                tx.Tx.Id,
		NumberOfAddress:   tx.Tx.NumberOfAddress,
		AccessAddressList: tx.Tx.AccessAddressList,
		Input:             tx.Tx.Input,
	})

	//cryp.Validate()

	if err != nil {
		log.Info(err)

	} else {
		//log.Info("PreCommitPhaseResponse All Receive., TxID: ", tx.Tx.Id, " Count: ", atomic.LoadUint64(&tx.ResponseCount))

		CommitPhase(tx) //启动Commit阶段
	}
	//log.Info("Done PreCommit")


}

