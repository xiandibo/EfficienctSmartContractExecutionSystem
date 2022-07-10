package server

import (
	"context"
	"errors"
	"github.com/vadiminshakov/committer/coordinatorserver"
	"github.com/vadiminshakov/committer/followerserver"
	"github.com/vadiminshakov/committer/helpers"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	status "google.golang.org/grpc/status"
	"net"
	"time"
)

// WhiteListChecker intercepts RPC and checks that the caller is whitelisted.
//捕获其他节点发过来的请求
func WhiteListChecker(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	peerinfo, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to retrieve peer info")
	}

	host, _, err := net.SplitHostPort(peerinfo.Addr.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	//serv := info.Server.(*Server)
	serv := info.Server.(*coordinatorserver.CoordinatorServer)
	if !helpers.Includes(serv.Config.Whitelist, host) {
		return nil, status.Errorf(codes.PermissionDenied, "host %s is not in whitelist", host)
	}

	// Calls the handler
	h, err := handler(ctx, req)

	return h, err
}

//无意义
func WhiteListCheckerForFollower(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	peerinfo, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to retrieve peer info")
	}

	host, _, err := net.SplitHostPort(peerinfo.Addr.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	//serv := info.Server.(*Server)
	serv := info.Server.(*followerserver.FollowerServer)
	if !helpers.Includes(serv.Config.Whitelist, host) {
		return nil, status.Errorf(codes.PermissionDenied, "host %s is not in whitelist", host)
	}

	// Calls the handler
	h, err := handler(ctx, req)

	return h, err
}

/*
  blocking interceptors for tests
*/

// PrecommitBlock blocks execution of all followers on precommit stage for 10s
func PrecommitBlockALL(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	if info.FullMethod == "/schema.Commit/Precommit" {
		time.Sleep(10000 * time.Millisecond)
	}

	// Calls the handler
	h, err := handler(ctx, req)

	return h, err
}

// PrecommitBlockCoordinator blocks execution of coordinator on precommit stage for 10s
func PrecommitBlockCoordinator(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	server, ok := info.Server.(*Server)
	if !ok {
		return nil, errors.New("failed to assert interface to Server type")
	}

	if server.Config.Role == "coordinator" {
		if info.FullMethod == "/schema.Commit/Precommit" {
			time.Sleep(10000 * time.Millisecond)
		}
	}

	// Calls the handler
	h, err := handler(ctx, req)

	return h, err
}
