// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.2
// source: model/kanjengMami.proto

package model

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

// CachingClient is the client API for Caching service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CachingClient interface {
	Put(ctx context.Context, in *CacheRequest, opts ...grpc.CallOption) (*Error, error)
	Get(ctx context.Context, in *CacheRequestKey, opts ...grpc.CallOption) (*CacheResponse, error)
	Delete(ctx context.Context, in *CacheRequestKey, opts ...grpc.CallOption) (*Error, error)
}

type cachingClient struct {
	cc grpc.ClientConnInterface
}

func NewCachingClient(cc grpc.ClientConnInterface) CachingClient {
	return &cachingClient{cc}
}

func (c *cachingClient) Put(ctx context.Context, in *CacheRequest, opts ...grpc.CallOption) (*Error, error) {
	out := new(Error)
	err := c.cc.Invoke(ctx, "/model.Caching/Put", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cachingClient) Get(ctx context.Context, in *CacheRequestKey, opts ...grpc.CallOption) (*CacheResponse, error) {
	out := new(CacheResponse)
	err := c.cc.Invoke(ctx, "/model.Caching/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cachingClient) Delete(ctx context.Context, in *CacheRequestKey, opts ...grpc.CallOption) (*Error, error) {
	out := new(Error)
	err := c.cc.Invoke(ctx, "/model.Caching/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CachingServer is the server API for Caching service.
// All implementations must embed UnimplementedCachingServer
// for forward compatibility
type CachingServer interface {
	Put(context.Context, *CacheRequest) (*Error, error)
	Get(context.Context, *CacheRequestKey) (*CacheResponse, error)
	Delete(context.Context, *CacheRequestKey) (*Error, error)
	mustEmbedUnimplementedCachingServer()
}

// UnimplementedCachingServer must be embedded to have forward compatible implementations.
type UnimplementedCachingServer struct {
}

func (UnimplementedCachingServer) Put(context.Context, *CacheRequest) (*Error, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Put not implemented")
}
func (UnimplementedCachingServer) Get(context.Context, *CacheRequestKey) (*CacheResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedCachingServer) Delete(context.Context, *CacheRequestKey) (*Error, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedCachingServer) mustEmbedUnimplementedCachingServer() {}

// UnsafeCachingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CachingServer will
// result in compilation errors.
type UnsafeCachingServer interface {
	mustEmbedUnimplementedCachingServer()
}

func RegisterCachingServer(s grpc.ServiceRegistrar, srv CachingServer) {
	s.RegisterService(&Caching_ServiceDesc, srv)
}

func _Caching_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CacheRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CachingServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/model.Caching/Put",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CachingServer).Put(ctx, req.(*CacheRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Caching_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CacheRequestKey)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CachingServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/model.Caching/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CachingServer).Get(ctx, req.(*CacheRequestKey))
	}
	return interceptor(ctx, in, info, handler)
}

func _Caching_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CacheRequestKey)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CachingServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/model.Caching/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CachingServer).Delete(ctx, req.(*CacheRequestKey))
	}
	return interceptor(ctx, in, info, handler)
}

// Caching_ServiceDesc is the grpc.ServiceDesc for Caching service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Caching_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "model.Caching",
	HandlerType: (*CachingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Put",
			Handler:    _Caching_Put_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _Caching_Get_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Caching_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "model/kanjengMami.proto",
}
