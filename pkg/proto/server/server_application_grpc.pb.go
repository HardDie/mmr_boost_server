// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: server_application.proto

package server

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

const (
	Application_Create_FullMethodName                     = "/gateway.Application/Create"
	Application_GetList_FullMethodName                    = "/gateway.Application/GetList"
	Application_GetItem_FullMethodName                    = "/gateway.Application/GetItem"
	Application_DeleteItem_FullMethodName                 = "/gateway.Application/DeleteItem"
	Application_GetManagementList_FullMethodName          = "/gateway.Application/GetManagementList"
	Application_GetManagementItem_FullMethodName          = "/gateway.Application/GetManagementItem"
	Application_UpdateManagementItem_FullMethodName       = "/gateway.Application/UpdateManagementItem"
	Application_GetManagementPrivateItem_FullMethodName   = "/gateway.Application/GetManagementPrivateItem"
	Application_UpdateManagementItemStatus_FullMethodName = "/gateway.Application/UpdateManagementItemStatus"
)

// ApplicationClient is the client API for Application service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ApplicationClient interface {
	// Create application for boosting
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
	// Getting a list of the applications you created
	GetList(ctx context.Context, in *GetListRequest, opts ...grpc.CallOption) (*GetListResponse, error)
	// Get the application you created
	GetItem(ctx context.Context, in *GetItemRequest, opts ...grpc.CallOption) (*GetItemResponse, error)
	// Delete created application
	DeleteItem(ctx context.Context, in *DeleteItemRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Getting a list of all applications. Access: admin, manager
	GetManagementList(ctx context.Context, in *GetManagementListRequest, opts ...grpc.CallOption) (*GetManagementListResponse, error)
	// Get the application by id
	GetManagementItem(ctx context.Context, in *GetManagementItemRequest, opts ...grpc.CallOption) (*GetManagementItemResponse, error)
	// Update application data
	UpdateManagementItem(ctx context.Context, in *UpdateManagementItemRequest, opts ...grpc.CallOption) (*UpdateManagementItemResponse, error)
	// Getting private information from an application by id
	GetManagementPrivateItem(ctx context.Context, in *GetManagementItemRequest, opts ...grpc.CallOption) (*GetManagementPrivateItemResponse, error)
	// Update application status
	UpdateManagementItemStatus(ctx context.Context, in *UpdateManagementItemStatusRequest, opts ...grpc.CallOption) (*UpdateManagementItemStatusResponse, error)
}

type applicationClient struct {
	cc grpc.ClientConnInterface
}

func NewApplicationClient(cc grpc.ClientConnInterface) ApplicationClient {
	return &applicationClient{cc}
}

func (c *applicationClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, Application_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *applicationClient) GetList(ctx context.Context, in *GetListRequest, opts ...grpc.CallOption) (*GetListResponse, error) {
	out := new(GetListResponse)
	err := c.cc.Invoke(ctx, Application_GetList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *applicationClient) GetItem(ctx context.Context, in *GetItemRequest, opts ...grpc.CallOption) (*GetItemResponse, error) {
	out := new(GetItemResponse)
	err := c.cc.Invoke(ctx, Application_GetItem_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *applicationClient) DeleteItem(ctx context.Context, in *DeleteItemRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Application_DeleteItem_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *applicationClient) GetManagementList(ctx context.Context, in *GetManagementListRequest, opts ...grpc.CallOption) (*GetManagementListResponse, error) {
	out := new(GetManagementListResponse)
	err := c.cc.Invoke(ctx, Application_GetManagementList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *applicationClient) GetManagementItem(ctx context.Context, in *GetManagementItemRequest, opts ...grpc.CallOption) (*GetManagementItemResponse, error) {
	out := new(GetManagementItemResponse)
	err := c.cc.Invoke(ctx, Application_GetManagementItem_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *applicationClient) UpdateManagementItem(ctx context.Context, in *UpdateManagementItemRequest, opts ...grpc.CallOption) (*UpdateManagementItemResponse, error) {
	out := new(UpdateManagementItemResponse)
	err := c.cc.Invoke(ctx, Application_UpdateManagementItem_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *applicationClient) GetManagementPrivateItem(ctx context.Context, in *GetManagementItemRequest, opts ...grpc.CallOption) (*GetManagementPrivateItemResponse, error) {
	out := new(GetManagementPrivateItemResponse)
	err := c.cc.Invoke(ctx, Application_GetManagementPrivateItem_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *applicationClient) UpdateManagementItemStatus(ctx context.Context, in *UpdateManagementItemStatusRequest, opts ...grpc.CallOption) (*UpdateManagementItemStatusResponse, error) {
	out := new(UpdateManagementItemStatusResponse)
	err := c.cc.Invoke(ctx, Application_UpdateManagementItemStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ApplicationServer is the server API for Application service.
// All implementations must embed UnimplementedApplicationServer
// for forward compatibility
type ApplicationServer interface {
	// Create application for boosting
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	// Getting a list of the applications you created
	GetList(context.Context, *GetListRequest) (*GetListResponse, error)
	// Get the application you created
	GetItem(context.Context, *GetItemRequest) (*GetItemResponse, error)
	// Delete created application
	DeleteItem(context.Context, *DeleteItemRequest) (*emptypb.Empty, error)
	// Getting a list of all applications. Access: admin, manager
	GetManagementList(context.Context, *GetManagementListRequest) (*GetManagementListResponse, error)
	// Get the application by id
	GetManagementItem(context.Context, *GetManagementItemRequest) (*GetManagementItemResponse, error)
	// Update application data
	UpdateManagementItem(context.Context, *UpdateManagementItemRequest) (*UpdateManagementItemResponse, error)
	// Getting private information from an application by id
	GetManagementPrivateItem(context.Context, *GetManagementItemRequest) (*GetManagementPrivateItemResponse, error)
	// Update application status
	UpdateManagementItemStatus(context.Context, *UpdateManagementItemStatusRequest) (*UpdateManagementItemStatusResponse, error)
	mustEmbedUnimplementedApplicationServer()
}

// UnimplementedApplicationServer must be embedded to have forward compatible implementations.
type UnimplementedApplicationServer struct {
}

func (UnimplementedApplicationServer) Create(context.Context, *CreateRequest) (*CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedApplicationServer) GetList(context.Context, *GetListRequest) (*GetListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetList not implemented")
}
func (UnimplementedApplicationServer) GetItem(context.Context, *GetItemRequest) (*GetItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetItem not implemented")
}
func (UnimplementedApplicationServer) DeleteItem(context.Context, *DeleteItemRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteItem not implemented")
}
func (UnimplementedApplicationServer) GetManagementList(context.Context, *GetManagementListRequest) (*GetManagementListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetManagementList not implemented")
}
func (UnimplementedApplicationServer) GetManagementItem(context.Context, *GetManagementItemRequest) (*GetManagementItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetManagementItem not implemented")
}
func (UnimplementedApplicationServer) UpdateManagementItem(context.Context, *UpdateManagementItemRequest) (*UpdateManagementItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateManagementItem not implemented")
}
func (UnimplementedApplicationServer) GetManagementPrivateItem(context.Context, *GetManagementItemRequest) (*GetManagementPrivateItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetManagementPrivateItem not implemented")
}
func (UnimplementedApplicationServer) UpdateManagementItemStatus(context.Context, *UpdateManagementItemStatusRequest) (*UpdateManagementItemStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateManagementItemStatus not implemented")
}
func (UnimplementedApplicationServer) mustEmbedUnimplementedApplicationServer() {}

// UnsafeApplicationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ApplicationServer will
// result in compilation errors.
type UnsafeApplicationServer interface {
	mustEmbedUnimplementedApplicationServer()
}

func RegisterApplicationServer(s grpc.ServiceRegistrar, srv ApplicationServer) {
	s.RegisterService(&Application_ServiceDesc, srv)
}

func _Application_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Application_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Application_GetList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServer).GetList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Application_GetList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServer).GetList(ctx, req.(*GetListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Application_GetItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServer).GetItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Application_GetItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServer).GetItem(ctx, req.(*GetItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Application_DeleteItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServer).DeleteItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Application_DeleteItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServer).DeleteItem(ctx, req.(*DeleteItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Application_GetManagementList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetManagementListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServer).GetManagementList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Application_GetManagementList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServer).GetManagementList(ctx, req.(*GetManagementListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Application_GetManagementItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetManagementItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServer).GetManagementItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Application_GetManagementItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServer).GetManagementItem(ctx, req.(*GetManagementItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Application_UpdateManagementItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateManagementItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServer).UpdateManagementItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Application_UpdateManagementItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServer).UpdateManagementItem(ctx, req.(*UpdateManagementItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Application_GetManagementPrivateItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetManagementItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServer).GetManagementPrivateItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Application_GetManagementPrivateItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServer).GetManagementPrivateItem(ctx, req.(*GetManagementItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Application_UpdateManagementItemStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateManagementItemStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationServer).UpdateManagementItemStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Application_UpdateManagementItemStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationServer).UpdateManagementItemStatus(ctx, req.(*UpdateManagementItemStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Application_ServiceDesc is the grpc.ServiceDesc for Application service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Application_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gateway.Application",
	HandlerType: (*ApplicationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Application_Create_Handler,
		},
		{
			MethodName: "GetList",
			Handler:    _Application_GetList_Handler,
		},
		{
			MethodName: "GetItem",
			Handler:    _Application_GetItem_Handler,
		},
		{
			MethodName: "DeleteItem",
			Handler:    _Application_DeleteItem_Handler,
		},
		{
			MethodName: "GetManagementList",
			Handler:    _Application_GetManagementList_Handler,
		},
		{
			MethodName: "GetManagementItem",
			Handler:    _Application_GetManagementItem_Handler,
		},
		{
			MethodName: "UpdateManagementItem",
			Handler:    _Application_UpdateManagementItem_Handler,
		},
		{
			MethodName: "GetManagementPrivateItem",
			Handler:    _Application_GetManagementPrivateItem_Handler,
		},
		{
			MethodName: "UpdateManagementItemStatus",
			Handler:    _Application_UpdateManagementItemStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "server_application.proto",
}
