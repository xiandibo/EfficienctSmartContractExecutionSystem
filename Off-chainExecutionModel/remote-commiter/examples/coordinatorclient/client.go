package main

import (
	"bytes"
	"context"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"github.com/vadiminshakov/committer/clientpeer"
	pb "github.com/vadiminshakov/committer/coordinator"
	"time"
)

const addr = "localhost:3000"
//const addr = "10.17.72.45:3000"
//const addr = "10.17.56.4:3000"

var kvstoreAddr = []string {
	"bd770416a3345f91e4b34576cb804a576fa48eb1",
	"5a443704dd4b594b382c22a083e2bd3090a6fef3",
}

var kvtestAddr = []string {
	"684c903c66d69777377f0945052160c9f778d689",
}

//原来的只涉及follower1和2的输入
var input1 = "da778c8a00000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000001a000000000000000000000000000000000000000000000000000000000000002a000000000000000000000000000000000000000000000000000000000000003a000000000000000000000000000000000000000000000000000000000000004a0000000000000000000000000bd770416a3345f91e4b34576cb804a576fa48eb10000000000000000000000000000000000000000000000000000000000000080000000000000000000000000bd770416a3345f91e4b34576cb804a576fa48eb100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000001310000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000013200000000000000000000000000000000000000000000000000000000000000000000000000000000000000bd770416a3345f91e4b34576cb804a576fa48eb10000000000000000000000000000000000000000000000000000000000000080000000000000000000000000bd770416a3345f91e4b34576cb804a576fa48eb100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000001310000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000013200000000000000000000000000000000000000000000000000000000000000000000000000000000000000bd770416a3345f91e4b34576cb804a576fa48eb10000000000000000000000000000000000000000000000000000000000000080000000000000000000000000bd770416a3345f91e4b34576cb804a576fa48eb100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000001310000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000013200000000000000000000000000000000000000000000000000000000000000000000000000000000000000bd770416a3345f91e4b34576cb804a576fa48eb10000000000000000000000000000000000000000000000000000000000000080000000000000000000000000bd770416a3345f91e4b34576cb804a576fa48eb100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000001310000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000013200000000000000000000000000000000000000000000000000000000000000000000000000000000000000bd770416a3345f91e4b34576cb804a576fa48eb10000000000000000000000000000000000000000000000000000000000000080000000000000000000000000bd770416a3345f91e4b34576cb804a576fa48eb100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000001310000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000013200000000000000000000000000000000000000000000000000000000000000"
//测试多几个follower的性能
var input2 = "da778c8a00000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000001a000000000000000000000000000000000000000000000000000000000000002a000000000000000000000000000000000000000000000000000000000000003a000000000000000000000000000000000000000000000000000000000000004a0000000000000000000000000bd770416a3345f91e4b34576cb804a576fa48eb100000000000000000000000000000000000000000000000000000000000000800000000000000000000000005a443704dd4b594b382c22a083e2bd3090a6fef300000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000023232000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000233330000000000000000000000000000000000000000000000000000000000000000000000000000000000005a443704dd4b594b382c22a083e2bd3090a6fef30000000000000000000000000000000000000000000000000000000000000080000000000000000000000000bd770416a3345f91e4b34576cb804a576fa48eb100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000002323200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000023333000000000000000000000000000000000000000000000000000000000000000000000000000000000000bd770416a3345f91e4b34576cb804a576fa48eb100000000000000000000000000000000000000000000000000000000000000800000000000000000000000005a443704dd4b594b382c22a083e2bd3090a6fef300000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000023232000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000235350000000000000000000000000000000000000000000000000000000000000000000000000000000000005a443704dd4b594b382c22a083e2bd3090a6fef30000000000000000000000000000000000000000000000000000000000000080000000000000000000000000bd770416a3345f91e4b34576cb804a576fa48eb100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000002323200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000023737000000000000000000000000000000000000000000000000000000000000000000000000000000000000bd770416a3345f91e4b34576cb804a576fa48eb100000000000000000000000000000000000000000000000000000000000000800000000000000000000000005a443704dd4b594b382c22a083e2bd3090a6fef300000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000002323200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000023939000000000000000000000000000000000000000000000000000000000000"


//启动客户端，向协调者发送交易
//客户端发送消息的入口
func main() {
	//tracer, err := trace.Tracer("client", addr)
	//if err != nil {
	//	panic(err)
	//}
	//Client API For commit service 建立与addr的连接
	//goweb3Test()

	cli, err := clientpeer.NewCoordinatorClientPeer(addr, nil)
	if err != nil {
		panic(err)
	}

	//此系统不支持并发
	//go concurrentTest(cli)

	//远程调用Put API
	//resp, err := cli.Put(context.Background(), "key3", []byte("1111"))

	//构造十个地址的accesslist
	AccessList := make([][]byte, 0)

	AccessList = append(AccessList, common.Hex2Bytes(kvtestAddr[0]))

	for i := 0; i < 10; i++ {
		AccessList = append(AccessList, common.Hex2Bytes(kvstoreAddr[0]))
	}

	var accessList []byte
	//log.Info(AccessList)
	for i := 0; i < 11; i++ {
		accessList = BytesCombine(AccessList...)
	}
	//包含11个地址的访问列表
	//log.Info(accessList)
	//log.Info("SendTransaction")
	j := 1
	for {
		for i := j * 500; i < (j + 1) * 500; i++ {
			_, err = cli.SendTransaction(
				context.Background(),
				&pb.Transaction{
					Id:                uint64(i),
					NumberOfAddress:   uint64(11),
					AccessAddressList: accessList,
					Input:             common.Hex2Bytes(input1),
				})
			if err != nil {
				panic(err)
			}
			log.Info("Send Transaction, AccessList: ", accessList)
			//fmt.Println(res.Result)
			//time.Sleep(10 * time.Millisecond)

		}
		log.Info("已发送：", (j+1)*100)
		time.Sleep(5 * time.Second)
		j++
	}
}

//BytesCombine 多个[]byte数组合并成一个[]byte
func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

// 测试go-web包
// 用于生成指定的对于合约的Input数据
//func goweb3Test()  {
//	typ := abi.MustNewType("tuple(address a, uint256 b)")
//
//	type Obj struct {
//		A web3.Address
//		B *big.Int
//	}
//	obj := &Obj{
//		A: web3.Address{0x1},
//		B: big.NewInt(1),
//	}
//
//	// Encode
//	encoded, err := typ.Encode(obj)
//	if err != nil {
//		panic(err)
//	}
//
//	// Decode output into a map
//	res, err := typ.Decode(encoded)
//	if err != nil {
//		panic(err)
//	}
//
//	// Decode into a struct
//	var obj2 Obj
//	if err := typ.DecodeStruct(encoded, &obj2); err != nil {
//		panic(err)
//	}
//
//	fmt.Println(res)
//	fmt.Println(obj)
//
//}