package addrtoip

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"time"
)

//此package存放智能合约地址对应的节点ip地址的映射
var (
	CoordinatorAddr = "10.20.82.234:3000" //协调者节点地址

	// AddressToIP map{合约地址:节点地址}，用于测试
	// 暂时一个IP代表一个共识组
	AddressToIP = map[common.Address]string{
		//kvstore
		common.HexToAddress("0xbd770416a3345f91e4b34576cb804a576fa48eb1"): "10.20.36.229:3001", //kvstore
		common.HexToAddress("0x5a443704dd4b594b382c22a083e2bd3090a6fef3"): "10.20.36.229:3003", //kvstore

		//kvtest
		common.HexToAddress("0x684c903c66d69777377f0945052160c9f778d689"): "10.20.36.229:3002", //kvtest
		common.HexToAddress("0x2ffea410b53de84d04c3b55e52a0d53fcfd0146e"): "10.20.36.229:3004", //kvtest

	}

	// NodeList 每个节点启动时，都按照这个列表初始化一次连接，这样不需要每次调用都建立连接，否则造成拒绝连接的情况
	NodeList = []string{
		"10.20.82.234:3000", // 协调者，不部署合约

		"10.20.3.46:3001",
		//"10.20.3.46:3002",
		//"10.20.3.46:3003",
		//"10.20.3.46:3004",
		//
		"10.20.6.199:3001",
		//"10.20.6.199:3002",
		//"10.20.6.199:3003",
		//"10.20.6.199:3004",
		//
		//"10.20.36.229:3001",
		//"10.20.36.229:3002",
		//"10.20.36.229:3003",
		//"10.20.36.229:3004",
		//
		//"10.20.0.57:3001",
		//"10.20.0.57:3002",
		//"10.20.0.57:3003",
		//"10.20.0.57:3004",
	}

	NodeAddrConnMap = map[string]*grpc.ClientConn{}


)

func EstablishConn() error {
	connParams := grpc.ConnectParams{
		Backoff: backoff.Config{
			BaseDelay: 100 * time.Millisecond,
			MaxDelay:  5 * time.Second,
		},
		MinConnectTimeout: 1000 * time.Millisecond,
	}
	for _, addr := range NodeList {
		conn, err := grpc.Dial(addr, grpc.WithConnectParams(connParams), grpc.WithInsecure())
		if err != nil {
			return errors.Wrap(err, "failed to connect")
		} else {
			NodeAddrConnMap[addr] = conn
		}

	}
	return nil
}