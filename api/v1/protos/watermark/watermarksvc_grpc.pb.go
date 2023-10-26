// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.1
// source: watermark/watermarksvc.proto

package watermark

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

const (
	Watermark_Create_FullMethodName        = "/watermark.Watermark/Create"
	Watermark_ServiceStatus_FullMethodName = "/watermark.Watermark/ServiceStatus"
)

// WatermarkClient is the client API for Watermark service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WatermarkClient interface {
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
	ServiceStatus(ctx context.Context, in *ServiceStatusRequest, opts ...grpc.CallOption) (*ServiceStatusResponse, error)
}

type watermarkClient struct {
	cc grpc.ClientConnInterface
}

func NewWatermarkClient(cc grpc.ClientConnInterface) WatermarkClient {
	return &watermarkClient{cc}
}

func (c *watermarkClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, Watermark_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *watermarkClient) ServiceStatus(ctx context.Context, in *ServiceStatusRequest, opts ...grpc.CallOption) (*ServiceStatusResponse, error) {
	out := new(ServiceStatusResponse)
	err := c.cc.Invoke(ctx, Watermark_ServiceStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WatermarkServer is the server API for Watermark service.
// All implementations must embed UnimplementedWatermarkServer
// for forward compatibility
type WatermarkServer interface {
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	ServiceStatus(context.Context, *ServiceStatusRequest) (*ServiceStatusResponse, error)
	mustEmbedUnimplementedWatermarkServer()
}

// UnimplementedWatermarkServer must be embedded to have forward compatible implementations.
type UnimplementedWatermarkServer struct {
}

func (UnimplementedWatermarkServer) Create(context.Context, *CreateRequest) (*CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedWatermarkServer) ServiceStatus(context.Context, *ServiceStatusRequest) (*ServiceStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ServiceStatus not implemented")
}
func (UnimplementedWatermarkServer) mustEmbedUnimplementedWatermarkServer() {}

// UnsafeWatermarkServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WatermarkServer will
// result in compilation errors.
type UnsafeWatermarkServer interface {
	mustEmbedUnimplementedWatermarkServer()
}

func RegisterWatermarkServer(s grpc.ServiceRegistrar, srv WatermarkServer) {
	s.RegisterService(&Watermark_ServiceDesc, srv)
}

func _Watermark_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WatermarkServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Watermark_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WatermarkServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Watermark_ServiceStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServiceStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WatermarkServer).ServiceStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Watermark_ServiceStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WatermarkServer).ServiceStatus(ctx, req.(*ServiceStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Watermark_ServiceDesc is the grpc.ServiceDesc for Watermark service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Watermark_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "watermark.Watermark",
	HandlerType: (*WatermarkServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Watermark_Create_Handler,
		},
		{
			MethodName: "ServiceStatus",
			Handler:    _Watermark_ServiceStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "watermark/watermarksvc.proto",
}
