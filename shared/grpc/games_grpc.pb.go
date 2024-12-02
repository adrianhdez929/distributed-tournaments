// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: grpc/games.proto

package grpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	GamesService_UploadGameFile_FullMethodName = "/games.GamesService/UploadGameFile"
)

// GamesServiceClient is the client API for GamesService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GamesServiceClient interface {
	UploadGameFile(ctx context.Context, in *UploadGameFileRequest, opts ...grpc.CallOption) (*UploadGameFileResponse, error)
}

type gamesServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGamesServiceClient(cc grpc.ClientConnInterface) GamesServiceClient {
	return &gamesServiceClient{cc}
}

func (c *gamesServiceClient) UploadGameFile(ctx context.Context, in *UploadGameFileRequest, opts ...grpc.CallOption) (*UploadGameFileResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UploadGameFileResponse)
	err := c.cc.Invoke(ctx, GamesService_UploadGameFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GamesServiceServer is the server API for GamesService service.
// All implementations must embed UnimplementedGamesServiceServer
// for forward compatibility.
type GamesServiceServer interface {
	UploadGameFile(context.Context, *UploadGameFileRequest) (*UploadGameFileResponse, error)
	mustEmbedUnimplementedGamesServiceServer()
}

// UnimplementedGamesServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedGamesServiceServer struct{}

func (UnimplementedGamesServiceServer) UploadGameFile(context.Context, *UploadGameFileRequest) (*UploadGameFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadGameFile not implemented")
}
func (UnimplementedGamesServiceServer) mustEmbedUnimplementedGamesServiceServer() {}
func (UnimplementedGamesServiceServer) testEmbeddedByValue()                      {}

// UnsafeGamesServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GamesServiceServer will
// result in compilation errors.
type UnsafeGamesServiceServer interface {
	mustEmbedUnimplementedGamesServiceServer()
}

func RegisterGamesServiceServer(s grpc.ServiceRegistrar, srv GamesServiceServer) {
	// If the following call pancis, it indicates UnimplementedGamesServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&GamesService_ServiceDesc, srv)
}

func _GamesService_UploadGameFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadGameFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GamesServiceServer).UploadGameFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GamesService_UploadGameFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GamesServiceServer).UploadGameFile(ctx, req.(*UploadGameFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GamesService_ServiceDesc is the grpc.ServiceDesc for GamesService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GamesService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "games.GamesService",
	HandlerType: (*GamesServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UploadGameFile",
			Handler:    _GamesService_UploadGameFile_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc/games.proto",
}
