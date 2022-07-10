// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package follower

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// FollowerClient is the client API for Follower service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FollowerClient interface {
	//rpc SC2PCPropose(ProposeReqForSC) returns (ResponseFromFollower); //coordinator rise propose to follower
	//rpc SC2PCCommit(CommitReqForSC) returns (ResponseFromFollower); //coordinatoe rise commit after propose phase
	InternalCallSmartContract(ctx context.Context, in *InputForSC, opts ...grpc.CallOption) (*ResultFromEVM, error)
	PreCommitRequest(ctx context.Context, in *Transaction, opts ...grpc.CallOption) (*emptypb.Empty, error)
	CommitRequest(ctx context.Context, in *Transaction, opts ...grpc.CallOption) (*emptypb.Empty, error)
	InternalCallSmartContractToVoter(ctx context.Context, in *InputForSC, opts ...grpc.CallOption) (*EVMResultFromVoter, error)
}

type followerClient struct {
	cc grpc.ClientConnInterface
}

func NewFollowerClient(cc grpc.ClientConnInterface) FollowerClient {
	return &followerClient{cc}
}

func (c *followerClient) InternalCallSmartContract(ctx context.Context, in *InputForSC, opts ...grpc.CallOption) (*ResultFromEVM, error) {
	out := new(ResultFromEVM)
	err := c.cc.Invoke(ctx, "/follower.Follower/InternalCallSmartContract", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *followerClient) PreCommitRequest(ctx context.Context, in *Transaction, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/follower.Follower/PreCommitRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *followerClient) CommitRequest(ctx context.Context, in *Transaction, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/follower.Follower/CommitRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *followerClient) InternalCallSmartContractToVoter(ctx context.Context, in *InputForSC, opts ...grpc.CallOption) (*EVMResultFromVoter, error) {
	out := new(EVMResultFromVoter)
	err := c.cc.Invoke(ctx, "/follower.Follower/InternalCallSmartContractToVoter", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FollowerServer is the server API for Follower service.
// All implementations must embed UnimplementedFollowerServer
// for forward compatibility
type FollowerServer interface {
	//rpc SC2PCPropose(ProposeReqForSC) returns (ResponseFromFollower); //coordinator rise propose to follower
	//rpc SC2PCCommit(CommitReqForSC) returns (ResponseFromFollower); //coordinatoe rise commit after propose phase
	InternalCallSmartContract(context.Context, *InputForSC) (*ResultFromEVM, error)
	PreCommitRequest(context.Context, *Transaction) (*emptypb.Empty, error)
	CommitRequest(context.Context, *Transaction) (*emptypb.Empty, error)
	InternalCallSmartContractToVoter(context.Context, *InputForSC) (*EVMResultFromVoter, error)
	mustEmbedUnimplementedFollowerServer()
}

// UnimplementedFollowerServer must be embedded to have forward compatible implementations.
type UnimplementedFollowerServer struct {
}

func (UnimplementedFollowerServer) InternalCallSmartContract(context.Context, *InputForSC) (*ResultFromEVM, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InternalCallSmartContract not implemented")
}
func (UnimplementedFollowerServer) PreCommitRequest(context.Context, *Transaction) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PreCommitRequest not implemented")
}
func (UnimplementedFollowerServer) CommitRequest(context.Context, *Transaction) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommitRequest not implemented")
}
func (UnimplementedFollowerServer) InternalCallSmartContractToVoter(context.Context, *InputForSC) (*EVMResultFromVoter, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InternalCallSmartContractToVoter not implemented")
}
func (UnimplementedFollowerServer) mustEmbedUnimplementedFollowerServer() {}

// UnsafeFollowerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FollowerServer will
// result in compilation errors.
type UnsafeFollowerServer interface {
	mustEmbedUnimplementedFollowerServer()
}

func RegisterFollowerServer(s grpc.ServiceRegistrar, srv FollowerServer) {
	s.RegisterService(&Follower_ServiceDesc, srv)
}

func _Follower_InternalCallSmartContract_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InputForSC)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FollowerServer).InternalCallSmartContract(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/follower.Follower/InternalCallSmartContract",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FollowerServer).InternalCallSmartContract(ctx, req.(*InputForSC))
	}
	return interceptor(ctx, in, info, handler)
}

func _Follower_PreCommitRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Transaction)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FollowerServer).PreCommitRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/follower.Follower/PreCommitRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FollowerServer).PreCommitRequest(ctx, req.(*Transaction))
	}
	return interceptor(ctx, in, info, handler)
}

func _Follower_CommitRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Transaction)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FollowerServer).CommitRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/follower.Follower/CommitRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FollowerServer).CommitRequest(ctx, req.(*Transaction))
	}
	return interceptor(ctx, in, info, handler)
}

func _Follower_InternalCallSmartContractToVoter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InputForSC)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FollowerServer).InternalCallSmartContractToVoter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/follower.Follower/InternalCallSmartContractToVoter",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FollowerServer).InternalCallSmartContractToVoter(ctx, req.(*InputForSC))
	}
	return interceptor(ctx, in, info, handler)
}

// Follower_ServiceDesc is the grpc.ServiceDesc for Follower service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Follower_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "follower.Follower",
	HandlerType: (*FollowerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "InternalCallSmartContract",
			Handler:    _Follower_InternalCallSmartContract_Handler,
		},
		{
			MethodName: "PreCommitRequest",
			Handler:    _Follower_PreCommitRequest_Handler,
		},
		{
			MethodName: "CommitRequest",
			Handler:    _Follower_CommitRequest_Handler,
		},
		{
			MethodName: "InternalCallSmartContractToVoter",
			Handler:    _Follower_InternalCallSmartContractToVoter_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "follower.proto",
}
