// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: nfs-api.proto

package api

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

// NFSSClient is the client API for NFSS service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NFSSClient interface {
	Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	Mount(ctx context.Context, in *MountRequest, opts ...grpc.CallOption) (*MountResponse, error)
	UnMount(ctx context.Context, in *UnMountRequest, opts ...grpc.CallOption) (*Empty, error)
	Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (NFSS_ReadClient, error)
	SaveAsync(ctx context.Context, opts ...grpc.CallOption) (NFSS_SaveAsyncClient, error)
	Save(ctx context.Context, in *SaveRequest, opts ...grpc.CallOption) (*Empty, error)
	Remove(ctx context.Context, in *RemoveRequest, opts ...grpc.CallOption) (*Empty, error)
	Chpem(ctx context.Context, in *ChpemRequest, opts ...grpc.CallOption) (*Empty, error)
}

type nFSSClient struct {
	cc grpc.ClientConnInterface
}

func NewNFSSClient(cc grpc.ClientConnInterface) NFSSClient {
	return &nFSSClient{cc}
}

func (c *nFSSClient) Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/api.NFSS/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nFSSClient) Mount(ctx context.Context, in *MountRequest, opts ...grpc.CallOption) (*MountResponse, error) {
	out := new(MountResponse)
	err := c.cc.Invoke(ctx, "/api.NFSS/Mount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nFSSClient) UnMount(ctx context.Context, in *UnMountRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/api.NFSS/UnMount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nFSSClient) Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (NFSS_ReadClient, error) {
	stream, err := c.cc.NewStream(ctx, &NFSS_ServiceDesc.Streams[0], "/api.NFSS/Read", opts...)
	if err != nil {
		return nil, err
	}
	x := &nFSSReadClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type NFSS_ReadClient interface {
	Recv() (*ReadResponse, error)
	grpc.ClientStream
}

type nFSSReadClient struct {
	grpc.ClientStream
}

func (x *nFSSReadClient) Recv() (*ReadResponse, error) {
	m := new(ReadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *nFSSClient) SaveAsync(ctx context.Context, opts ...grpc.CallOption) (NFSS_SaveAsyncClient, error) {
	stream, err := c.cc.NewStream(ctx, &NFSS_ServiceDesc.Streams[1], "/api.NFSS/SaveAsync", opts...)
	if err != nil {
		return nil, err
	}
	x := &nFSSSaveAsyncClient{stream}
	return x, nil
}

type NFSS_SaveAsyncClient interface {
	Send(*SaveRequest) error
	CloseAndRecv() (*Empty, error)
	grpc.ClientStream
}

type nFSSSaveAsyncClient struct {
	grpc.ClientStream
}

func (x *nFSSSaveAsyncClient) Send(m *SaveRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *nFSSSaveAsyncClient) CloseAndRecv() (*Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *nFSSClient) Save(ctx context.Context, in *SaveRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/api.NFSS/Save", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nFSSClient) Remove(ctx context.Context, in *RemoveRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/api.NFSS/Remove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nFSSClient) Chpem(ctx context.Context, in *ChpemRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/api.NFSS/Chpem", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NFSSServer is the server API for NFSS service.
// All implementations should embed UnimplementedNFSSServer
// for forward compatibility
type NFSSServer interface {
	Ping(context.Context, *Empty) (*Empty, error)
	Mount(context.Context, *MountRequest) (*MountResponse, error)
	UnMount(context.Context, *UnMountRequest) (*Empty, error)
	Read(*ReadRequest, NFSS_ReadServer) error
	SaveAsync(NFSS_SaveAsyncServer) error
	Save(context.Context, *SaveRequest) (*Empty, error)
	Remove(context.Context, *RemoveRequest) (*Empty, error)
	Chpem(context.Context, *ChpemRequest) (*Empty, error)
}

// UnimplementedNFSSServer should be embedded to have forward compatible implementations.
type UnimplementedNFSSServer struct {
}

func (UnimplementedNFSSServer) Ping(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedNFSSServer) Mount(context.Context, *MountRequest) (*MountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Mount not implemented")
}
func (UnimplementedNFSSServer) UnMount(context.Context, *UnMountRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnMount not implemented")
}
func (UnimplementedNFSSServer) Read(*ReadRequest, NFSS_ReadServer) error {
	return status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (UnimplementedNFSSServer) SaveAsync(NFSS_SaveAsyncServer) error {
	return status.Errorf(codes.Unimplemented, "method SaveAsync not implemented")
}
func (UnimplementedNFSSServer) Save(context.Context, *SaveRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Save not implemented")
}
func (UnimplementedNFSSServer) Remove(context.Context, *RemoveRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
}
func (UnimplementedNFSSServer) Chpem(context.Context, *ChpemRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Chpem not implemented")
}

// UnsafeNFSSServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NFSSServer will
// result in compilation errors.
type UnsafeNFSSServer interface {
	mustEmbedUnimplementedNFSSServer()
}

func RegisterNFSSServer(s grpc.ServiceRegistrar, srv NFSSServer) {
	s.RegisterService(&NFSS_ServiceDesc, srv)
}

func _NFSS_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NFSSServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.NFSS/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NFSSServer).Ping(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _NFSS_Mount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NFSSServer).Mount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.NFSS/Mount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NFSSServer).Mount(ctx, req.(*MountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NFSS_UnMount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnMountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NFSSServer).UnMount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.NFSS/UnMount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NFSSServer).UnMount(ctx, req.(*UnMountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NFSS_Read_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ReadRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(NFSSServer).Read(m, &nFSSReadServer{stream})
}

type NFSS_ReadServer interface {
	Send(*ReadResponse) error
	grpc.ServerStream
}

type nFSSReadServer struct {
	grpc.ServerStream
}

func (x *nFSSReadServer) Send(m *ReadResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _NFSS_SaveAsync_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(NFSSServer).SaveAsync(&nFSSSaveAsyncServer{stream})
}

type NFSS_SaveAsyncServer interface {
	SendAndClose(*Empty) error
	Recv() (*SaveRequest, error)
	grpc.ServerStream
}

type nFSSSaveAsyncServer struct {
	grpc.ServerStream
}

func (x *nFSSSaveAsyncServer) SendAndClose(m *Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *nFSSSaveAsyncServer) Recv() (*SaveRequest, error) {
	m := new(SaveRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _NFSS_Save_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NFSSServer).Save(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.NFSS/Save",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NFSSServer).Save(ctx, req.(*SaveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NFSS_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NFSSServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.NFSS/Remove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NFSSServer).Remove(ctx, req.(*RemoveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NFSS_Chpem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChpemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NFSSServer).Chpem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.NFSS/Chpem",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NFSSServer).Chpem(ctx, req.(*ChpemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NFSS_ServiceDesc is the grpc.ServiceDesc for NFSS service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NFSS_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.NFSS",
	HandlerType: (*NFSSServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _NFSS_Ping_Handler,
		},
		{
			MethodName: "Mount",
			Handler:    _NFSS_Mount_Handler,
		},
		{
			MethodName: "UnMount",
			Handler:    _NFSS_UnMount_Handler,
		},
		{
			MethodName: "Save",
			Handler:    _NFSS_Save_Handler,
		},
		{
			MethodName: "Remove",
			Handler:    _NFSS_Remove_Handler,
		},
		{
			MethodName: "Chpem",
			Handler:    _NFSS_Chpem_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Read",
			Handler:       _NFSS_Read_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SaveAsync",
			Handler:       _NFSS_SaveAsync_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "nfs-api.proto",
}
