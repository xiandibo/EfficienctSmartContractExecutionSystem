package main
//
//import (
//	"context"
//	"fmt"
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/openzipkin/zipkin-go"
//	log "github.com/sirupsen/logrus"
//	"github.com/status-im/keycard-go/hexutils"
//	"github.com/stretchr/testify/assert"
//	"github.com/umbracle/go-web3"
//	"github.com/umbracle/go-web3/abi"
//	"github.com/vadiminshakov/committer/config"
//	"github.com/vadiminshakov/committer/hooks"
//	"github.com/vadiminshakov/committer/peer"
//	pb "github.com/vadiminshakov/committer/proto"
//	"github.com/vadiminshakov/committer/server"
//	"github.com/vadiminshakov/committer/trace"
//	"google.golang.org/grpc"
//	"os"
//	"strconv"
//	"testing"
//	"time"
//)
//
//const (
//	COORDINATOR_TYPE = "coordinator"
//	FOLLOWER_TYPE    = "follower"
//	BADGER_DIR       = "/tmp/badger"
//)
//
//const (
//	NOT_BLOCKING = iota
//	BLOCK_ON_PRECOMMIT_FOLLOWERS
//	BLOCK_ON_PRECOMMIT_COORDINATOR
//)
//
//var (
//	whitelist = []string{"127.0.0.1"}
//	nodes     = map[string][]*config.Config{
//		COORDINATOR_TYPE: {
//			{Nodeaddr: "localhost:3000", Role: "coordinator",
//				Followers: []string{"localhost:3001", "localhost:3002", "localhost:3003", "localhost:3004", "localhost:3005"},
//				Whitelist: whitelist, CommitType: "two-phase", Timeout: 1000, WithTrace: false},
//			{Nodeaddr: "localhost:5000", Role: "coordinator",
//				Followers: []string{"localhost:3001", "localhost:3002", "localhost:3003", "localhost:3004", "localhost:3005"},
//				Whitelist: whitelist, CommitType: "three-phase", Timeout: 1000, WithTrace: false},
//		},
//		FOLLOWER_TYPE: {
//			&config.Config{Nodeaddr: "localhost:3001", Role: "follower", Coordinator: "localhost:3000", Whitelist: whitelist, Timeout: 1000, WithTrace: false},
//			&config.Config{Nodeaddr: "localhost:3002", Role: "follower", Coordinator: "localhost:3000", Whitelist: whitelist, Timeout: 1000, WithTrace: false},
//			&config.Config{Nodeaddr: "localhost:3003", Role: "follower", Coordinator: "localhost:3000", Whitelist: whitelist, Timeout: 1000, WithTrace: false},
//			&config.Config{Nodeaddr: "localhost:3004", Role: "follower", Coordinator: "localhost:3000", Whitelist: whitelist, Timeout: 1000, WithTrace: false},
//			&config.Config{Nodeaddr: "localhost:3005", Role: "follower", Coordinator: "localhost:3000", Whitelist: whitelist, Timeout: 1000, WithTrace: false},
//		},
//	}
//)
//
//var testtable = map[string][]byte{
//	"key1": []byte("value1"),
//	"key2": []byte("value2"),
//	"key3": []byte("value3"),
//	"key4": []byte("value4"),
//	"key5": []byte("value5"),
//	"key6": []byte("value6"),
//	"key7": []byte("value7"),
//	"key8": []byte("value8"),
//}
//
//func TestHappyPath(t *testing.T) {
//	log.SetLevel(log.FatalLevel)
//
//	done := make(chan struct{})
//	go startnodes(NOT_BLOCKING, done)
//	time.Sleep(6 * time.Second) // wait for coordinators and followers to start and establish connections
//
//	var height uint64 = 0
//	for _, coordConfig := range nodes[COORDINATOR_TYPE] {
//		if coordConfig.CommitType == "two-phase" {
//			log.Println("***\nTEST IN TWO-PHASE MODE\n***")
//		} else {
//			log.Println("***\nTEST IN THREE-PHASE MODE\n***")
//		}
//		var (
//			tracer *zipkin.Tracer
//			err    error
//		)
//		if coordConfig.WithTrace {
//			tracer, err = trace.Tracer("client", coordConfig.Nodeaddr)
//			if err != nil {
//				t.Errorf("no tracer, err: %v", err)
//			}
//		}
//		c, err := peer.New(coordConfig.Nodeaddr, tracer)
//		if err != nil {
//			t.Error(err)
//		}
//
//		for key, val := range testtable {
//			resp, err := c.Put(context.Background(), key, val)
//			if err != nil {
//				t.Error(err)
//			}
//			if resp.Type != pb.Type_ACK {
//				t.Error("msg is not acknowledged")
//			}
//			// ok, value is added, let's increment height counter
//			height++
//		}
//
//		// connect to followers and check that them added key-value
//		for _, node := range nodes[FOLLOWER_TYPE] {
//			if coordConfig.WithTrace {
//				tracer, err = trace.Tracer(fmt.Sprintf("%s:%s", coordConfig.Role, coordConfig.Nodeaddr), coordConfig.Nodeaddr)
//				if err != nil {
//					t.Errorf("no tracer, err: %v", err)
//				}
//			}
//			cli, err := peer.New(node.Nodeaddr, tracer)
//			assert.NoError(t, err, "err not nil")
//			for key, val := range testtable {
//				// check values added by nodes
//				resp, err := cli.Get(context.Background(), key)
//				assert.NoError(t, err, "err not nil")
//				assert.Equal(t, resp.Value, val)
//
//				// check height of node
//				nodeInfo, err := cli.NodeInfo(context.Background())
//				assert.NoError(t, err, "err not nil")
//				assert.Equal(t, nodeInfo.Height, height, "node %s ahead, %d commits behind (current height is %d)", node.Nodeaddr, height-nodeInfo.Height, nodeInfo.Height)
//			}
//		}
//	}
//	done <- struct{}{}
//	time.Sleep(2 * time.Second)
//}
//
//// 5 followers, 1 coordinator
//// on precommit stage all followers stops responding
//func Test_3PC_6NODES_ALLFAILURE_ON_PRECOMMIT(t *testing.T) {
//	log.SetLevel(log.FatalLevel)
//
//	done := make(chan struct{})
//	go startnodes(BLOCK_ON_PRECOMMIT_FOLLOWERS, done)
//	time.Sleep(10 * time.Second) // wait for coordinators and followers to start and establish connections
//
//	var (
//		tracer *zipkin.Tracer
//		err    error
//	)
//	if nodes[COORDINATOR_TYPE][1].WithTrace {
//		tracer, err = trace.Tracer("client", nodes[COORDINATOR_TYPE][1].Nodeaddr)
//		if err != nil {
//			t.Errorf("no tracer, err: %v", err)
//		}
//	}
//	c, err := peer.New(nodes[COORDINATOR_TYPE][1].Nodeaddr, tracer)
//	assert.NoError(t, err, "err not nil")
//	for key, val := range testtable {
//		resp, err := c.Put(context.Background(), key, val)
//		assert.NoError(t, err, "err not nil")
//		assert.NotEqual(t, resp.Type, pb.Type_ACK, "msg shouldn't be acknowledged")
//	}
//
//	done <- struct{}{}
//	time.Sleep(2 * time.Second)
//}
//
//// 5 followers, 1 coordinator
//// on precommit stage coordinator stops responding
//func Test_3PC_6NODES_COORDINATORFAILURE_ON_PRECOMMIT(t *testing.T) {
//	log.SetLevel(log.FatalLevel)
//	done := make(chan struct{})
//	go startnodes(BLOCK_ON_PRECOMMIT_COORDINATOR, done)
//	time.Sleep(10 * time.Second) // wait for coordinators and followers to start and establish connections
//
//	var (
//		tracer *zipkin.Tracer
//		err    error
//	)
//	if nodes[COORDINATOR_TYPE][1].WithTrace {
//		tracer, err = trace.Tracer("client", nodes[COORDINATOR_TYPE][1].Nodeaddr)
//		if err != nil {
//			t.Errorf("no tracer, err: %v", err)
//		}
//	}
//
//	c, err := peer.New(nodes[COORDINATOR_TYPE][1].Nodeaddr, tracer)
//	assert.NoError(t, err, "err not nil")
//
//	var height uint64 = 0
//	for key, val := range testtable {
//		resp, err := c.Put(context.Background(), key, val)
//		assert.NoError(t, err, "err not nil")
//		assert.Equal(t, resp.Type, pb.Type_ACK, "msg should be acknowledged")
//		height += 1
//	}
//
//	// connect to follower and check that them added key-value
//	for _, node := range nodes[FOLLOWER_TYPE] {
//		cli, err := peer.New(node.Nodeaddr, tracer)
//		assert.NoError(t, err, "err not nil")
//		for key, val := range testtable {
//			// check values added by nodes
//			resp, err := cli.Get(context.Background(), key)
//			assert.NoError(t, err, "err not nil")
//			assert.Equal(t, resp.Value, val)
//
//			// check height of node
//			nodeInfo, err := cli.NodeInfo(context.Background())
//			assert.NoError(t, err, "err not nil")
//			assert.Equal(t, nodeInfo.Height, height, "node %s ahead, %d commits behind (current height is %d)", node.Nodeaddr, height-nodeInfo.Height, nodeInfo.Height)
//		}
//	}
//
//	done <- struct{}{}
//	time.Sleep(2 * time.Second)
//}
//
//func startnodes(block int, done chan struct{}) {
//	COORDINATOR_BADGER := fmt.Sprintf("%s%s%d", BADGER_DIR, "coordinator", time.Now().UnixNano())
//	FOLLOWER_BADGER := fmt.Sprintf("%s%s%d", BADGER_DIR, "follower", time.Now().UnixNano())
//
//	os.Mkdir(COORDINATOR_BADGER, os.FileMode(0777))
//	os.Mkdir(FOLLOWER_BADGER, os.FileMode(0777))
//
//	var blocking grpc.UnaryServerInterceptor
//	switch block {
//	case BLOCK_ON_PRECOMMIT_FOLLOWERS:
//		blocking = server.PrecommitBlockALL
//	case BLOCK_ON_PRECOMMIT_COORDINATOR:
//		blocking = server.PrecommitBlockCoordinator
//	}
//
//	// start followers
//	for i, node := range nodes[FOLLOWER_TYPE] {
//		// create db dir
//		os.Mkdir(node.DBPath, os.FileMode(0777))
//		node.DBPath = fmt.Sprintf("%s%s%s", FOLLOWER_BADGER, strconv.Itoa(i), "~")
//		// start follower
//		hooks, err := hooks.Get()
//		if err != nil {
//			panic(err)
//		}
//		followerServer, err := server.NewCommitServer(node, hooks...)
//		if err != nil {
//			panic(err)
//		}
//
//		if block == 0 {
//			go followerServer.Run(server.WhiteListChecker)
//		} else {
//			go followerServer.Run(server.WhiteListChecker, blocking)
//		}
//
//		defer followerServer.Stop()
//	}
//	time.Sleep(3 * time.Second)
//
//	// start coordinators (in two- and three-phase modes)
//	for i, coordConfig := range nodes[COORDINATOR_TYPE] {
//		// create db dir
//		os.Mkdir(coordConfig.DBPath, os.FileMode(0777))
//		coordConfig.DBPath = fmt.Sprintf("%s%s%s", COORDINATOR_BADGER, strconv.Itoa(i), "~")
//		// start coordinator
//		hooks, err := hooks.Get()
//		if err != nil {
//			panic(err)
//		}
//		coordServer, err := server.NewCommitServer(coordConfig, hooks...)
//		if err != nil {
//			panic(err)
//		}
//
//		if block == 0 {
//			go coordServer.Run(server.WhiteListChecker)
//		} else {
//			go coordServer.Run(server.WhiteListChecker, blocking)
//		}
//
//		defer coordServer.Stop()
//	}
//
//	<-done
//	// prune
//	os.RemoveAll(BADGER_DIR)
//}
//
//func Test_test(t *testing.T) {
//	typ := abi.MustNewType("tuple(address A11, string B12, address C13, string D14, " +
//		"address A21, string B22, address C23, string D24, " +
//		"address A31, string B32, address C33, string D34, " +
//		"address A41, string B42, address C43, string D44, " +
//		"address A51, string B52, address C53, string D54)")
//
//	type Obj struct {
//		A11 web3.Address
//		B12 string
//		C13 web3.Address
//		D14 string
//		A21 web3.Address
//		B22 string
//		C23 web3.Address
//		D24 string
//		A31 web3.Address
//		B32 string
//		C33 web3.Address
//		D34 string
//		A41 web3.Address
//		B42 string
//		C43 web3.Address
//		D44 string
//		A51 web3.Address
//		B52 string
//		C53 web3.Address
//		D54 string
//	}
//	obj := &Obj{
//		A11: web3.HexToAddress("0xbd770416a3345f91e4b34576cb804a576fa48eb1"),
//		B12: "1",
//		C13: web3.HexToAddress("0xbd770416a3345f91e4b34576cb804a576fa48eb1"),
//		D14: "2",
//		A21: web3.HexToAddress("0xbd770416a3345f91e4b34576cb804a576fa48eb1"),
//		B22: "1",
//		C23: web3.HexToAddress("0xbd770416a3345f91e4b34576cb804a576fa48eb1"),
//		D24: "2",
//		A31: web3.HexToAddress("0xbd770416a3345f91e4b34576cb804a576fa48eb1"),
//		B32: "1",
//		C33: web3.HexToAddress("0xbd770416a3345f91e4b34576cb804a576fa48eb1"),
//		D34: "2",
//		A41: web3.HexToAddress("0xbd770416a3345f91e4b34576cb804a576fa48eb1"),
//		B42: "1",
//		C43: web3.HexToAddress("0xbd770416a3345f91e4b34576cb804a576fa48eb1"),
//		D44: "2",
//		A51: web3.HexToAddress("0xbd770416a3345f91e4b34576cb804a576fa48eb1"),
//		B52: "1",
//		C53: web3.HexToAddress("0xbd770416a3345f91e4b34576cb804a576fa48eb1"),
//		D54: "2",
//	}
//
//	objSlice := make([]*Obj, 0)
//	objSlice = append(objSlice, obj)
//	objSlice = append(objSlice, obj)
//
//
//	// Encode
//	encoded, err := typ.Encode(obj)
//	//encoded, err := encode(typ, obj)
//	if err != nil {
//		panic(err)
//	}
//
//	encodedtest := make([]byte, 0)
//	encodedtest = append(encodedtest, hexutils.HexToBytes("da778c8a")...)
//
//
//
//		encodedtemp, err := typ.Encode(obj)
//		if err != nil {
//			panic(err)
//		}
//		encodedtest = append(encodedtest, encodedtemp...)
//
//
//	fmt.Printf("%X\n", encodedtest)
//	var input1 = "da778c8a00000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000001a000000000000000000000000000000000000000000000000000000000000002a000000000000000000000000000000000000000000000000000000000000003a000000000000000000000000000000000000000000000000000000000000004a0000000000000000000000000bd770416a3345f91e4b34576cb804a576fa48eb10000000000000000000000000000000000000000000000000000000000000080000000000000000000000000bd770416a3345f91e4b34576cb804a576fa48eb100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000001310000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000013200000000000000000000000000000000000000000000000000000000000000000000000000000000000000bd770416a3345f91e4b34576cb804a576fa48eb10000000000000000000000000000000000000000000000000000000000000080000000000000000000000000bd770416a3345f91e4b34576cb804a576fa48eb100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000001310000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000013200000000000000000000000000000000000000000000000000000000000000000000000000000000000000bd770416a3345f91e4b34576cb804a576fa48eb10000000000000000000000000000000000000000000000000000000000000080000000000000000000000000bd770416a3345f91e4b34576cb804a576fa48eb100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000001310000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000013200000000000000000000000000000000000000000000000000000000000000000000000000000000000000bd770416a3345f91e4b34576cb804a576fa48eb10000000000000000000000000000000000000000000000000000000000000080000000000000000000000000bd770416a3345f91e4b34576cb804a576fa48eb100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000001310000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000013200000000000000000000000000000000000000000000000000000000000000000000000000000000000000bd770416a3345f91e4b34576cb804a576fa48eb10000000000000000000000000000000000000000000000000000000000000080000000000000000000000000bd770416a3345f91e4b34576cb804a576fa48eb100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000001310000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000013200000000000000000000000000000000000000000000000000000000000000"
//	fmt.Printf("%X\n", common.Hex2Bytes(input1))
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
//	fmt.Println(encoded)
//	fmt.Println(res)
//	fmt.Println(obj)
//}
//
//func encode(t *abi.Type, args ...interface{}) ([]byte, error) {
//	return abi.Encode(args, t)
//	//typ.Encode(args)
//}