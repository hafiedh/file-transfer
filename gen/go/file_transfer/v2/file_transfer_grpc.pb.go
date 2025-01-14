// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: file_transfer/v2/file_transfer.proto

package v2

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
	FileTransferService_UploadFile_FullMethodName      = "/file_transfer.v2.FileTransferService/UploadFile"
	FileTransferService_GetPresignedURL_FullMethodName = "/file_transfer.v2.FileTransferService/GetPresignedURL"
)

// FileTransferServiceClient is the client API for FileTransferService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FileTransferServiceClient interface {
	UploadFile(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[UploadFileRequest, UploadFileResponse], error)
	GetPresignedURL(ctx context.Context, in *GetPresignedURLRequest, opts ...grpc.CallOption) (*GetPresignedURLResponse, error)
}

type fileTransferServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFileTransferServiceClient(cc grpc.ClientConnInterface) FileTransferServiceClient {
	return &fileTransferServiceClient{cc}
}

func (c *fileTransferServiceClient) UploadFile(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[UploadFileRequest, UploadFileResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &FileTransferService_ServiceDesc.Streams[0], FileTransferService_UploadFile_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[UploadFileRequest, UploadFileResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type FileTransferService_UploadFileClient = grpc.ClientStreamingClient[UploadFileRequest, UploadFileResponse]

func (c *fileTransferServiceClient) GetPresignedURL(ctx context.Context, in *GetPresignedURLRequest, opts ...grpc.CallOption) (*GetPresignedURLResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPresignedURLResponse)
	err := c.cc.Invoke(ctx, FileTransferService_GetPresignedURL_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FileTransferServiceServer is the server API for FileTransferService service.
// All implementations must embed UnimplementedFileTransferServiceServer
// for forward compatibility.
type FileTransferServiceServer interface {
	UploadFile(grpc.ClientStreamingServer[UploadFileRequest, UploadFileResponse]) error
	GetPresignedURL(context.Context, *GetPresignedURLRequest) (*GetPresignedURLResponse, error)
	mustEmbedUnimplementedFileTransferServiceServer()
}

// UnimplementedFileTransferServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedFileTransferServiceServer struct{}

func (UnimplementedFileTransferServiceServer) UploadFile(grpc.ClientStreamingServer[UploadFileRequest, UploadFileResponse]) error {
	return status.Errorf(codes.Unimplemented, "method UploadFile not implemented")
}
func (UnimplementedFileTransferServiceServer) GetPresignedURL(context.Context, *GetPresignedURLRequest) (*GetPresignedURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPresignedURL not implemented")
}
func (UnimplementedFileTransferServiceServer) mustEmbedUnimplementedFileTransferServiceServer() {}
func (UnimplementedFileTransferServiceServer) testEmbeddedByValue()                             {}

// UnsafeFileTransferServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FileTransferServiceServer will
// result in compilation errors.
type UnsafeFileTransferServiceServer interface {
	mustEmbedUnimplementedFileTransferServiceServer()
}

func RegisterFileTransferServiceServer(s grpc.ServiceRegistrar, srv FileTransferServiceServer) {
	// If the following call pancis, it indicates UnimplementedFileTransferServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&FileTransferService_ServiceDesc, srv)
}

func _FileTransferService_UploadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FileTransferServiceServer).UploadFile(&grpc.GenericServerStream[UploadFileRequest, UploadFileResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type FileTransferService_UploadFileServer = grpc.ClientStreamingServer[UploadFileRequest, UploadFileResponse]

func _FileTransferService_GetPresignedURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPresignedURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileTransferServiceServer).GetPresignedURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FileTransferService_GetPresignedURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileTransferServiceServer).GetPresignedURL(ctx, req.(*GetPresignedURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FileTransferService_ServiceDesc is the grpc.ServiceDesc for FileTransferService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FileTransferService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "file_transfer.v2.FileTransferService",
	HandlerType: (*FileTransferServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPresignedURL",
			Handler:    _FileTransferService_GetPresignedURL_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadFile",
			Handler:       _FileTransferService_UploadFile_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "file_transfer/v2/file_transfer.proto",
}
