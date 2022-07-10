package clientpeer

import (
	"context"
	"github.com/openzipkin/zipkin-go"
	zipkingrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	"github.com/pkg/errors"
	"github.com/vadiminshakov/committer/addrtoip"
	pb "github.com/vadiminshakov/committer/follower"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

//调用远程节点的接口
type FollowerClient struct {
	Connection pb.FollowerClient
	Tracer     *zipkin.Tracer
}

// NewFollowerClientPeer creates instance of peer client.
// 'addr' is a coordinator network address (host + port).
func NewFollowerClientPeer(addr string, tracer *zipkin.Tracer) (*FollowerClient, error) {
	if _, ok := addrtoip.NodeAddrConnMap[addr]; ok {
		conn := addrtoip.NodeAddrConnMap[addr]
		return &FollowerClient{Connection: pb.NewFollowerClient(conn), Tracer: tracer}, nil
	} else {
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
		return &FollowerClient{Connection: pb.NewFollowerClient(conn), Tracer: tracer}, nil
	}


	/*var (
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
	return &FollowerClient{Connection: pb.NewFollowerClient(conn), Tracer: tracer}, nil*/
}

func (client *FollowerClient) InternalCallSmartContract(ctx context.Context, in *pb.InputForSC) (*pb.ResultFromEVM, error) {
	return client.Connection.InternalCallSmartContract(ctx, in)
}

func (client *FollowerClient) InternalCallSmartContractToVoter(ctx context.Context, in *pb.InputForSC) (*pb.EVMResultFromVoter, error) {
	return client.Connection.InternalCallSmartContractToVoter(ctx, in)
}

func (client *FollowerClient) PreCommitRequest(ctx context.Context, in *pb.Transaction) (*emptypb.Empty, error) {
	return client.Connection.PreCommitRequest(ctx, in)
}

func (client *FollowerClient) CommitRequest(ctx context.Context, in *pb.Transaction) (*emptypb.Empty, error) {
	return client.Connection.CommitRequest(ctx, in)
}


