package txpool

import (
	"github.com/ethereum/go-ethereum/common"
	"testing"
)

func TestTxpool_IsSerializable(t *testing.T) {
	txpool := NewTxpool()

	tx := Tx{
		Tx:            nil,
		AccessList:    nil,
		State:         false,
		ResponseCount: 0,
	}
	tx.AccessList = append(tx.AccessList, common.HexToAddress("0xbd770416a3345f91e4b34576cb804a576fa48eb1"))
	tx.AccessList = append(tx.AccessList, common.HexToAddress("0x5a443704dd4b594b382c22a083e2bd3090a6fef3"))
	tx.AccessList = append(tx.AccessList, common.HexToAddress("0x684c903c66d69777377f0945052160c9f778d689"))
	tx.AccessList = append(tx.AccessList, common.HexToAddress("0x2ffea410b53de84d04c3b55e52a0d53fcfd0146e"))
	tx.AccessList = append(tx.AccessList, common.HexToAddress("0xbd770416a3345f91e4b34576cb804a576fa48eb1"))


	txpool.IsSerializable(&tx)

}
