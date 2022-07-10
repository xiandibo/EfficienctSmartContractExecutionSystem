package clientpeer

import (
	"context"
	"github.com/openzipkin/zipkin-go"
	zipkingrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	"github.com/pkg/errors"
	"github.com/vadiminshakov/committer/addrtoip"
	pb "github.com/vadiminshakov/committer/coordinator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

//调用远程节点的接口
type CoordinatorClient struct {
	Connection pb.CoordinatorClient
	Tracer     *zipkin.Tracer
}

// New creates instance of peer client.
// 'addr' is a coordinator network address (host + port).
func NewCoordinatorClientPeer(addr string, tracer *zipkin.Tracer) (*CoordinatorClient, error) {
	if _, ok := addrtoip.NodeAddrConnMap[addr]; ok {
		conn := addrtoip.NodeAddrConnMap[addr]
		return &CoordinatorClient{Connection: pb.NewCoordinatorClient(conn), Tracer: tracer}, nil

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
		return &CoordinatorClient{Connection: pb.NewCoordinatorClient(conn), Tracer: tracer}, nil
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
	return &CoordinatorClient{Connection: pb.NewCoordinatorClient(conn), Tracer: tracer}, nil*/
}

func (client *CoordinatorClient) SendTransaction(ctx context.Context, in *pb.Transaction) (*emptypb.Empty, error) {
	return client.Connection.SendTransaction(ctx, in)
}

func (client *CoordinatorClient) CommitPhaseResponse(ctx context.Context, in *pb.ACK) (*emptypb.Empty, error) {
	return client.Connection.CommitPhaseResponse(ctx, in)
}

func (client *CoordinatorClient) PreCommitPhaseResponse(ctx context.Context, in *pb.ACK) (*emptypb.Empty, error) {
	return client.Connection.PreCommitPhaseResponse(ctx, in)
}


