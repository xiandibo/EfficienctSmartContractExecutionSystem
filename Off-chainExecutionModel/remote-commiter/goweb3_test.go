package main
//
//import (
//	"fmt"
//	"github.com/umbracle/go-web3"
//	"github.com/umbracle/go-web3/abi"
//	"math/big"
//	"testing"
//)
//
//func test(t *testing.T) {
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
//}