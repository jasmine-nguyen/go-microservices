// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.1
// source: proto/money_movement_svc.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	MoneyMovementService_Authorize_FullMethodName = "/MoneyMovementService/Authorize"
	MoneyMovementService_Capture_FullMethodName   = "/MoneyMovementService/Capture"
)

// MoneyMovementServiceClient is the client API for MoneyMovementService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MoneyMovementServiceClient interface {
	Authorize(ctx context.Context, in *AuthorizeRequest, opts ...grpc.CallOption) (*AuthorizeResponse, error)
	Capture(ctx context.Context, in *CaptureRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type moneyMovementServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMoneyMovementServiceClient(cc grpc.ClientConnInterface) MoneyMovementServiceClient {
	return &moneyMovementServiceClient{cc}
}

func (c *moneyMovementServiceClient) Authorize(ctx context.Context, in *AuthorizeRequest, opts ...grpc.CallOption) (*AuthorizeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AuthorizeResponse)
	err := c.cc.Invoke(ctx, MoneyMovementService_Authorize_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *moneyMovementServiceClient) Capture(ctx context.Context, in *CaptureRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, MoneyMovementService_Capture_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MoneyMovementServiceServer is the server API for MoneyMovementService service.
// All implementations must embed UnimplementedMoneyMovementServiceServer
// for forward compatibility.
type MoneyMovementServiceServer interface {
	Authorize(context.Context, *AuthorizeRequest) (*AuthorizeResponse, error)
	Capture(context.Context, *CaptureRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedMoneyMovementServiceServer()
}

// UnimplementedMoneyMovementServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMoneyMovementServiceServer struct{}

func (UnimplementedMoneyMovementServiceServer) Authorize(context.Context, *AuthorizeRequest) (*AuthorizeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authorize not implemented")
}
func (UnimplementedMoneyMovementServiceServer) Capture(context.Context, *CaptureRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Capture not implemented")
}
func (UnimplementedMoneyMovementServiceServer) mustEmbedUnimplementedMoneyMovementServiceServer() {}
func (UnimplementedMoneyMovementServiceServer) testEmbeddedByValue()                              {}

// UnsafeMoneyMovementServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MoneyMovementServiceServer will
// result in compilation errors.
type UnsafeMoneyMovementServiceServer interface {
	mustEmbedUnimplementedMoneyMovementServiceServer()
}

func RegisterMoneyMovementServiceServer(s grpc.ServiceRegistrar, srv MoneyMovementServiceServer) {
	// If the following call pancis, it indicates UnimplementedMoneyMovementServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&MoneyMovementService_ServiceDesc, srv)
}

func _MoneyMovementService_Authorize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthorizeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MoneyMovementServiceServer).Authorize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MoneyMovementService_Authorize_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MoneyMovementServiceServer).Authorize(ctx, req.(*AuthorizeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MoneyMovementService_Capture_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CaptureRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MoneyMovementServiceServer).Capture(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MoneyMovementService_Capture_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MoneyMovementServiceServer).Capture(ctx, req.(*CaptureRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MoneyMovementService_ServiceDesc is the grpc.ServiceDesc for MoneyMovementService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MoneyMovementService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "MoneyMovementService",
	HandlerType: (*MoneyMovementServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Authorize",
			Handler:    _MoneyMovementService_Authorize_Handler,
		},
		{
			MethodName: "Capture",
			Handler:    _MoneyMovementService_Capture_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/money_movement_svc.proto",
}