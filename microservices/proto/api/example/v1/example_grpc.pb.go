// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.6
// source: example/v1/example.proto

package examplepb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ExampleServiceClient is the client API for ExampleService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExampleServiceClient interface {
	Hello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error)
	World(ctx context.Context, in *WorldRequest, opts ...grpc.CallOption) (*WorldResponse, error)
	StreamWorld(ctx context.Context, opts ...grpc.CallOption) (ExampleService_StreamWorldClient, error)
}

type exampleServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewExampleServiceClient(cc grpc.ClientConnInterface) ExampleServiceClient {
	return &exampleServiceClient{cc}
}

func (c *exampleServiceClient) Hello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error) {
	out := new(HelloResponse)
	err := c.cc.Invoke(ctx, "/example.v1.ExampleService/Hello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exampleServiceClient) World(ctx context.Context, in *WorldRequest, opts ...grpc.CallOption) (*WorldResponse, error) {
	out := new(WorldResponse)
	err := c.cc.Invoke(ctx, "/example.v1.ExampleService/World", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exampleServiceClient) StreamWorld(ctx context.Context, opts ...grpc.CallOption) (ExampleService_StreamWorldClient, error) {
	stream, err := c.cc.NewStream(ctx, &ExampleService_ServiceDesc.Streams[0], "/example.v1.ExampleService/StreamWorld", opts...)
	if err != nil {
		return nil, err
	}
	x := &exampleServiceStreamWorldClient{stream}
	return x, nil
}

type ExampleService_StreamWorldClient interface {
	Send(*StreamWorldRequest) error
	Recv() (*StreamWorldResponse, error)
	grpc.ClientStream
}

type exampleServiceStreamWorldClient struct {
	grpc.ClientStream
}

func (x *exampleServiceStreamWorldClient) Send(m *StreamWorldRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *exampleServiceStreamWorldClient) Recv() (*StreamWorldResponse, error) {
	m := new(StreamWorldResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ExampleServiceServer is the server API for ExampleService service.
// All implementations must embed UnimplementedExampleServiceServer
// for forward compatibility
type ExampleServiceServer interface {
	Hello(context.Context, *HelloRequest) (*HelloResponse, error)
	World(context.Context, *WorldRequest) (*WorldResponse, error)
	StreamWorld(ExampleService_StreamWorldServer) error
	mustEmbedUnimplementedExampleServiceServer()
}

// UnimplementedExampleServiceServer must be embedded to have forward compatible implementations.
type UnimplementedExampleServiceServer struct {
}

func (UnimplementedExampleServiceServer) Hello(context.Context, *HelloRequest) (*HelloResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Hello not implemented")
}
func (UnimplementedExampleServiceServer) World(context.Context, *WorldRequest) (*WorldResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method World not implemented")
}
func (UnimplementedExampleServiceServer) StreamWorld(ExampleService_StreamWorldServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamWorld not implemented")
}
func (UnimplementedExampleServiceServer) mustEmbedUnimplementedExampleServiceServer() {}

// UnsafeExampleServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExampleServiceServer will
// result in compilation errors.
type UnsafeExampleServiceServer interface {
	mustEmbedUnimplementedExampleServiceServer()
}

func RegisterExampleServiceServer(s grpc.ServiceRegistrar, srv ExampleServiceServer) {
	s.RegisterService(&ExampleService_ServiceDesc, srv)
}

func _ExampleService_Hello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExampleServiceServer).Hello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/example.v1.ExampleService/Hello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExampleServiceServer).Hello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExampleService_World_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WorldRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExampleServiceServer).World(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/example.v1.ExampleService/World",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExampleServiceServer).World(ctx, req.(*WorldRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExampleService_StreamWorld_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ExampleServiceServer).StreamWorld(&exampleServiceStreamWorldServer{stream})
}

type ExampleService_StreamWorldServer interface {
	Send(*StreamWorldResponse) error
	Recv() (*StreamWorldRequest, error)
	grpc.ServerStream
}

type exampleServiceStreamWorldServer struct {
	grpc.ServerStream
}

func (x *exampleServiceStreamWorldServer) Send(m *StreamWorldResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *exampleServiceStreamWorldServer) Recv() (*StreamWorldRequest, error) {
	m := new(StreamWorldRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ExampleService_ServiceDesc is the grpc.ServiceDesc for ExampleService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ExampleService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "example.v1.ExampleService",
	HandlerType: (*ExampleServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Hello",
			Handler:    _ExampleService_Hello_Handler,
		},
		{
			MethodName: "World",
			Handler:    _ExampleService_World_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamWorld",
			Handler:       _ExampleService_StreamWorld_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "example/v1/example.proto",
}
