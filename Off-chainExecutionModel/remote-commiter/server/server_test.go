package server

import (
	"bytes"
	"context"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	pb "github.com/vadiminshakov/committer/proto"
	"testing"
)

var kvstoreAddr = []string {
	"bd770416a3345f91e4b34576cb804a576fa48eb1",
	"5a443704dd4b594b382c22a083e2bd3090a6fef3",
}

var kvtestAddr = []string {
	"684c903c66d69777377f0945052160c9f778d689",
}


func TestSendTx(t *testing.T) {

	//NodeCache := cache.New()
	//NodeCache.Set(1, "testkey", []byte("testvalue"))

	//构造十个地址的accesslist
	AccessList := make([][]byte, 0)

	AccessList = append(AccessList, common.Hex2Bytes(kvtestAddr[0]))

	for i := 0; i < 10; i++ {
		AccessList = append(AccessList, common.Hex2Bytes(kvstoreAddr[0]))
	}

	var accessList []byte
	log.Info(AccessList)
	for i := 0; i < 11; i++ {
		accessList = BytesCombine(AccessList...)
	}
	//包含11个地址的访问列表
	log.Info(accessList)

	s := &Server{}
	tx := &pb.Transaction{
		Id: uint64(1),
		NumberOfAddress: uint64(11),
		AccessAddressList:accessList,
	}
	s.SendTx(context.Background(), tx)



}

//BytesCombine 多个[]byte数组合并成一个[]byte
func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}
