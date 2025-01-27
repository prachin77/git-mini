// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.0--rc1
// source: background_service.proto

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
	BackgroundService_InitWorkspaceConnWithPort_FullMethodName = "/BackgroundService/InitWorkspaceConnWithPort"
	BackgroundService_GetFiles_FullMethodName                  = "/BackgroundService/GetFiles"
	BackgroundService_GetHostPcPublicKey_FullMethodName        = "/BackgroundService/GetHostPcPublicKey"
)

// BackgroundServiceClient is the client API for BackgroundService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BackgroundServiceClient interface {
	// Initialiize new folder/workspace into config file
	InitWorkspaceConnWithPort(ctx context.Context, in *InitRequest, opts ...grpc.CallOption) (*InitResponse, error)
	GetFiles(ctx context.Context, in *CloneRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Files], error)
	// Gets Host PC's (from which we're cloning files) public Keys
	GetHostPcPublicKey(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PublicKey, error)
}

type backgroundServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBackgroundServiceClient(cc grpc.ClientConnInterface) BackgroundServiceClient {
	return &backgroundServiceClient{cc}
}

func (c *backgroundServiceClient) InitWorkspaceConnWithPort(ctx context.Context, in *InitRequest, opts ...grpc.CallOption) (*InitResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(InitResponse)
	err := c.cc.Invoke(ctx, BackgroundService_InitWorkspaceConnWithPort_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backgroundServiceClient) GetFiles(ctx context.Context, in *CloneRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Files], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &BackgroundService_ServiceDesc.Streams[0], BackgroundService_GetFiles_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[CloneRequest, Files]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type BackgroundService_GetFilesClient = grpc.ServerStreamingClient[Files]

func (c *backgroundServiceClient) GetHostPcPublicKey(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PublicKey, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PublicKey)
	err := c.cc.Invoke(ctx, BackgroundService_GetHostPcPublicKey_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BackgroundServiceServer is the server API for BackgroundService service.
// All implementations must embed UnimplementedBackgroundServiceServer
// for forward compatibility.
type BackgroundServiceServer interface {
	// Initialiize new folder/workspace into config file
	InitWorkspaceConnWithPort(context.Context, *InitRequest) (*InitResponse, error)
	GetFiles(*CloneRequest, grpc.ServerStreamingServer[Files]) error
	// Gets Host PC's (from which we're cloning files) public Keys
	GetHostPcPublicKey(context.Context, *emptypb.Empty) (*PublicKey, error)
	mustEmbedUnimplementedBackgroundServiceServer()
}

// UnimplementedBackgroundServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedBackgroundServiceServer struct{}

func (UnimplementedBackgroundServiceServer) InitWorkspaceConnWithPort(context.Context, *InitRequest) (*InitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitWorkspaceConnWithPort not implemented")
}
func (UnimplementedBackgroundServiceServer) GetFiles(*CloneRequest, grpc.ServerStreamingServer[Files]) error {
	return status.Errorf(codes.Unimplemented, "method GetFiles not implemented")
}
func (UnimplementedBackgroundServiceServer) GetHostPcPublicKey(context.Context, *emptypb.Empty) (*PublicKey, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHostPcPublicKey not implemented")
}
func (UnimplementedBackgroundServiceServer) mustEmbedUnimplementedBackgroundServiceServer() {}
func (UnimplementedBackgroundServiceServer) testEmbeddedByValue()                           {}

// UnsafeBackgroundServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BackgroundServiceServer will
// result in compilation errors.
type UnsafeBackgroundServiceServer interface {
	mustEmbedUnimplementedBackgroundServiceServer()
}

func RegisterBackgroundServiceServer(s grpc.ServiceRegistrar, srv BackgroundServiceServer) {
	// If the following call pancis, it indicates UnimplementedBackgroundServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&BackgroundService_ServiceDesc, srv)
}

func _BackgroundService_InitWorkspaceConnWithPort_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackgroundServiceServer).InitWorkspaceConnWithPort(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BackgroundService_InitWorkspaceConnWithPort_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackgroundServiceServer).InitWorkspaceConnWithPort(ctx, req.(*InitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BackgroundService_GetFiles_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(CloneRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BackgroundServiceServer).GetFiles(m, &grpc.GenericServerStream[CloneRequest, Files]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type BackgroundService_GetFilesServer = grpc.ServerStreamingServer[Files]

func _BackgroundService_GetHostPcPublicKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackgroundServiceServer).GetHostPcPublicKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BackgroundService_GetHostPcPublicKey_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackgroundServiceServer).GetHostPcPublicKey(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// BackgroundService_ServiceDesc is the grpc.ServiceDesc for BackgroundService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BackgroundService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "BackgroundService",
	HandlerType: (*BackgroundServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "InitWorkspaceConnWithPort",
			Handler:    _BackgroundService_InitWorkspaceConnWithPort_Handler,
		},
		{
			MethodName: "GetHostPcPublicKey",
			Handler:    _BackgroundService_GetHostPcPublicKey_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetFiles",
			Handler:       _BackgroundService_GetFiles_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "background_service.proto",
}
