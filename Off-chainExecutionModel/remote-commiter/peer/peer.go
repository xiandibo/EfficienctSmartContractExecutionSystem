package peer

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/openzipkin/zipkin-go"
	zipkingrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	"github.com/pkg/errors"
	pb "github.com/vadiminshakov/committer/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"time"
)


//调用远程节点的接口
type CommitClient struct {
	Connection pb.CommitClient
	Tracer     *zipkin.Tracer
}

// New creates instance of peer client.
// 'addr' is a coordinator network address (host + port).
func New(addr string, tracer *zipkin.Tracer) (*CommitClient, error) {
	var (
		conn *grpc.ClientConn
		err  error
	)

	connParams := grpc.ConnectParams{
		Backoff: backoff.Config{
			BaseDelay: 100 * time.Millisecond,
			MaxDelay:  5 * time.Second,
		},
		MinConnectTimeout: 200 * time.Millisecond,
	}
	if tracer != nil {
		conn, err = grpc.Dial(addr, grpc.WithConnectParams(connParams), grpc.WithInsecure(), grpc.WithStatsHandler(zipkingrpc.NewClientHandler(tracer)))
	} else {
		conn, err = grpc.Dial(addr, grpc.WithConnectParams(connParams), grpc.WithInsecure())
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect")
	}
	return &CommitClient{Connection: pb.NewCommitClient(conn), Tracer: tracer}, nil
}

func (client *CommitClient) Propose(ctx context.Context, req *pb.ProposeRequest) (*pb.Response, error) {
	var span zipkin.Span
	if client.Tracer != nil {
		span, ctx = client.Tracer.StartSpanFromContext(ctx, "Propose")
		defer span.Finish()
	}
	return client.Connection.Propose(ctx, req)
}

func (client *CommitClient) Precommit(ctx context.Context, req *pb.PrecommitRequest) (*pb.Response, error) {
	var span zipkin.Span
	if client.Tracer != nil {
		span, ctx = client.Tracer.StartSpanFromContext(ctx, "Precommit")
		defer span.Finish()
	}
	return client.Connection.Precommit(ctx, req)
}

func (client *CommitClient) Commit(ctx context.Context, req *pb.CommitRequest) (*pb.Response, error) {
	var span zipkin.Span
	if client.Tracer != nil {
		span, ctx = client.Tracer.StartSpanFromContext(ctx, "Commit")
		defer span.Finish()
	}
	return client.Connection.Commit(ctx, req)
}

// Put sends key/value pair to peer (it should be a coordinator).
// The coordinator reaches consensus and all peers commit the value.
func (client *CommitClient) Put(ctx context.Context, key string, value []byte) (*pb.Response, error) {
	var span zipkin.Span
	if client.Tracer != nil {
		span, ctx = client.Tracer.StartSpanFromContext(ctx, "Put")
		defer span.Finish()
	}
	return client.Connection.Put(ctx, &pb.Entry{Key: key, Value: value})
}

// NodeInfo gets info about current node height
func (client *CommitClient) NodeInfo(ctx context.Context) (*pb.Info, error) {
	var span zipkin.Span
	if client.Tracer != nil {
		span, ctx = client.Tracer.StartSpanFromContext(ctx, "NodeInfo")
		defer span.Finish()
	}
	return client.Connection.NodeInfo(ctx, &empty.Empty{})
}

// Get queries value of specific key
func (client *CommitClient) Get(ctx context.Context, key string) (*pb.Value, error) {
	var span zipkin.Span
	if client.Tracer != nil {
		span, ctx = client.Tracer.StartSpanFromContext(ctx, "Get")
		defer span.Finish()
	}
	return client.Connection.Get(ctx, &pb.Msg{Key: key})
}

// CallSmartContract 调用智能合约，相当于客户端发起一笔调用智能合约的事务
// 此函数被客户端本地调用
// 此处调用远程节点方法
// 后续应作为客户端调用智能合约的入口
func (client *CommitClient) CallSmartContract(ctx context.Context, input *pb.InputForSC) (*pb.ResultFromEVM, error) {
	var span zipkin.Span
	if client.Tracer != nil {
		span, ctx = client.Tracer.StartSpanFromContext(ctx, "Get")
		defer span.Finish()
	}

	fmt.Println("客户端本地发起事务，调用协调者方法 == CallSmartContract")
	return client.Connection.CallSmartContract(ctx, input)
}