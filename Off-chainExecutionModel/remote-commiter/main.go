package main

import (
	"bufio"
	"fmt"
	"github.com/vadiminshakov/committer/addrtoip"
	"github.com/vadiminshakov/committer/config"
	"github.com/vadiminshakov/committer/coordinatorserver"
	"github.com/vadiminshakov/committer/followerserver"
	"github.com/vadiminshakov/committer/goroutinepool"
	"github.com/vadiminshakov/committer/scexecution"
	"github.com/vadiminshakov/committer/txpool"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)


func main() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//获取启动命令行参数
	conf := config.Get()
	//proposehook or commithook
	//hooks, err := hooks.Get()
	//if err != nil {
	//	panic(err)
	//}
	//有可能是coordinator 或 follower
	//以后需要编写不同的server来区分不同节点的方法，全部放在commitServer太过凌乱
	//s, err := server.NewCommitServer(conf, hooks...)
	if conf.Role == "coordinator" {
		//coordinatorserver
		s, err := coordinatorserver.NewCoordinatorServer(conf)

		if err != nil {
			panic(err)
		}

		//初始化合约代码和状态，初始化设置
		if s.ConfigOfChain == nil {
			s.ConfigOfChain = new(scexecution.Config)
			scexecution.Init(s.ConfigOfChain)
			s.Txpool.InitAddrMutexMap()
		}
		addrtoip.EstablishConn()

		//s.Run(server.WhiteListChecker)
		s.Run()

		signChan := make(chan struct{})
		go getCommandArgs(signChan)
		go runWorker(signChan, s)

		<-ch
		s.ExitCh <- *new(struct{})
		s.Stop()

	} else if conf.Role == "follower" {
		runtime.GOMAXPROCS(2)
		//followerserver
		s, err := followerserver.NewFollowerServer(conf)

		if err != nil {
			panic(err)
		}

		//初始化合约代码和状态，初始化设置
		if s.ConfigOfChain == nil {
			s.ConfigOfChain = new(scexecution.Config)
			scexecution.Init(s.ConfigOfChain)
			s.Txpool.InitAddrMutexMap()
		}
		addrtoip.EstablishConn()
		addrtoip.EstablishConnWithVoter(s.Addr)


		//s.Run(server.WhiteListCheckerForFollower)
		s.Run()
		<-ch
		s.ExitFlag = true
		s.Stop()

	}



}

// 获取用户输入的命令参数
// 目前只能输入一次
func getCommandArgs(ch chan struct{}) {
	buf := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	sentence, err := buf.ReadBytes('\n')
	if err != nil {
		fmt.Println(err)
	}
	cmd := string(sentence)
	fmt.Println(cmd)
	if cmd == "run\n" {
		// 启动执行
		ch <- struct{}{}
	} else {
		fmt.Println("非法的启动命令")
	}
}

// 监听ch通道，如果接收到启动信号，就启动执行tx
func runWorker(ch chan struct{}, s *coordinatorserver.CoordinatorServer)  {
	// 阻塞等待启动信号
	<-ch

	//创建一个协程池,最大开启3个协程worker
	p := goroutinepool.NewPool(40000)

	m := s.Txpool.TxIn2PCCMap.Items()
	for k, v := range m {
		fmt.Println(k)
		//创建一个Task
		ta := goroutinepool.NewTask(func(in interface{}) error {
			//fmt.Println(time.Now())
			s.Arise2PCPreCommitRequest(in.(*txpool.Tx))
			return nil
		},
		v)

		p.EntryChannel <- ta
	}

	fmt.Println("EntryChannel:", len(p.EntryChannel))
	//go watcher(p)
	// 用于计算用时
	coordinatorserver.Start = time.Now()
	coordinatorserver.TotalTx = uint64(len(p.EntryChannel))
	p.Run()
	// fmt.Println("EntryChannel:", len(p.EntryChannel))

}

func watcher(p *goroutinepool.Pool)  {
	start := time.Now()
	for {
		if len(p.JobsChannel) == 0 && len(p.EntryChannel) == 0 {
			end := time.Now()
			fmt.Println("全部执行完成！用时：", end.Sub(start).Seconds())
			break
		}
	}

}