// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: SupportService.proto

package events

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

// SupportClient is the client API for Support service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SupportClient interface {
	GetNotifications(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Notifies, error)
	SetNotified(ctx context.Context, in *NotificationIDReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
	CleanupOldEvents(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type supportClient struct {
	cc grpc.ClientConnInterface
}

func NewSupportClient(cc grpc.ClientConnInterface) SupportClient {
	return &supportClient{cc}
}

func (c *supportClient) GetNotifications(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Notifies, error) {
	out := new(Notifies)
	err := c.cc.Invoke(ctx, "/api.support/GetNotifications", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *supportClient) SetNotified(ctx context.Context, in *NotificationIDReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/api.support/SetNotified", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *supportClient) CleanupOldEvents(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/api.support/CleanupOldEvents", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SupportServer is the server API for Support service.
// All implementations must embed UnimplementedSupportServer
// for forward compatibility
type SupportServer interface {
	GetNotifications(context.Context, *emptypb.Empty) (*Notifies, error)
	SetNotified(context.Context, *NotificationIDReq) (*emptypb.Empty, error)
	CleanupOldEvents(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	mustEmbedUnimplementedSupportServer()
}

// UnimplementedSupportServer must be embedded to have forward compatible implementations.
type UnimplementedSupportServer struct {
}

func (UnimplementedSupportServer) GetNotifications(context.Context, *emptypb.Empty) (*Notifies, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNotifications not implemented")
}
func (UnimplementedSupportServer) SetNotified(context.Context, *NotificationIDReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetNotified not implemented")
}
func (UnimplementedSupportServer) CleanupOldEvents(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CleanupOldEvents not implemented")
}
func (UnimplementedSupportServer) mustEmbedUnimplementedSupportServer() {}

// UnsafeSupportServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SupportServer will
// result in compilation errors.
type UnsafeSupportServer interface {
	mustEmbedUnimplementedSupportServer()
}

func RegisterSupportServer(s grpc.ServiceRegistrar, srv SupportServer) {
	s.RegisterService(&Support_ServiceDesc, srv)
}

func _Support_GetNotifications_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SupportServer).GetNotifications(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.support/GetNotifications",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SupportServer).GetNotifications(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Support_SetNotified_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotificationIDReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SupportServer).SetNotified(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.support/SetNotified",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SupportServer).SetNotified(ctx, req.(*NotificationIDReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Support_CleanupOldEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SupportServer).CleanupOldEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.support/CleanupOldEvents",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SupportServer).CleanupOldEvents(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Support_ServiceDesc is the grpc.ServiceDesc for Support service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Support_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.support",
	HandlerType: (*SupportServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetNotifications",
			Handler:    _Support_GetNotifications_Handler,
		},
		{
			MethodName: "SetNotified",
			Handler:    _Support_SetNotified_Handler,
		},
		{
			MethodName: "CleanupOldEvents",
			Handler:    _Support_CleanupOldEvents_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "SupportService.proto",
}
