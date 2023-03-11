// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.22.0
// source: loms.proto

package loms

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

// LOMSClient is the client API for LOMS service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LOMSClient interface {
	// CreateOrder creates a new order for a user, reserves ordered products in a warehouse.
	CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error)
	// ListOrder lists order information.
	ListOrder(ctx context.Context, in *ListOrderRequest, opts ...grpc.CallOption) (*ListOrderResponse, error)
	// OrderPaid marks order as paid.
	OrderPaid(ctx context.Context, in *OrderPaidRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// CancelOrder cancels order, makes previously reserved products available.
	CancelOrder(ctx context.Context, in *CancelOrderRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Stocks returns a number of available products with a given SKU in different warehouses.
	Stocks(ctx context.Context, in *StocksRequest, opts ...grpc.CallOption) (*StocksResponse, error)
}

type lOMSClient struct {
	cc grpc.ClientConnInterface
}

func NewLOMSClient(cc grpc.ClientConnInterface) LOMSClient {
	return &lOMSClient{cc}
}

func (c *lOMSClient) CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error) {
	out := new(CreateOrderResponse)
	err := c.cc.Invoke(ctx, "/loms.LOMS/CreateOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lOMSClient) ListOrder(ctx context.Context, in *ListOrderRequest, opts ...grpc.CallOption) (*ListOrderResponse, error) {
	out := new(ListOrderResponse)
	err := c.cc.Invoke(ctx, "/loms.LOMS/ListOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lOMSClient) OrderPaid(ctx context.Context, in *OrderPaidRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/loms.LOMS/OrderPaid", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lOMSClient) CancelOrder(ctx context.Context, in *CancelOrderRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/loms.LOMS/CancelOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lOMSClient) Stocks(ctx context.Context, in *StocksRequest, opts ...grpc.CallOption) (*StocksResponse, error) {
	out := new(StocksResponse)
	err := c.cc.Invoke(ctx, "/loms.LOMS/Stocks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LOMSServer is the server API for LOMS service.
// All implementations must embed UnimplementedLOMSServer
// for forward compatibility
type LOMSServer interface {
	// CreateOrder creates a new order for a user, reserves ordered products in a warehouse.
	CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error)
	// ListOrder lists order information.
	ListOrder(context.Context, *ListOrderRequest) (*ListOrderResponse, error)
	// OrderPaid marks order as paid.
	OrderPaid(context.Context, *OrderPaidRequest) (*emptypb.Empty, error)
	// CancelOrder cancels order, makes previously reserved products available.
	CancelOrder(context.Context, *CancelOrderRequest) (*emptypb.Empty, error)
	// Stocks returns a number of available products with a given SKU in different warehouses.
	Stocks(context.Context, *StocksRequest) (*StocksResponse, error)
	mustEmbedUnimplementedLOMSServer()
}

// UnimplementedLOMSServer must be embedded to have forward compatible implementations.
type UnimplementedLOMSServer struct {
}

func (UnimplementedLOMSServer) CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrder not implemented")
}
func (UnimplementedLOMSServer) ListOrder(context.Context, *ListOrderRequest) (*ListOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListOrder not implemented")
}
func (UnimplementedLOMSServer) OrderPaid(context.Context, *OrderPaidRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrderPaid not implemented")
}
func (UnimplementedLOMSServer) CancelOrder(context.Context, *CancelOrderRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelOrder not implemented")
}
func (UnimplementedLOMSServer) Stocks(context.Context, *StocksRequest) (*StocksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stocks not implemented")
}
func (UnimplementedLOMSServer) mustEmbedUnimplementedLOMSServer() {}

// UnsafeLOMSServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LOMSServer will
// result in compilation errors.
type UnsafeLOMSServer interface {
	mustEmbedUnimplementedLOMSServer()
}

func RegisterLOMSServer(s grpc.ServiceRegistrar, srv LOMSServer) {
	s.RegisterService(&LOMS_ServiceDesc, srv)
}

func _LOMS_CreateOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LOMSServer).CreateOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loms.LOMS/CreateOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LOMSServer).CreateOrder(ctx, req.(*CreateOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LOMS_ListOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LOMSServer).ListOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loms.LOMS/ListOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LOMSServer).ListOrder(ctx, req.(*ListOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LOMS_OrderPaid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderPaidRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LOMSServer).OrderPaid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loms.LOMS/OrderPaid",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LOMSServer).OrderPaid(ctx, req.(*OrderPaidRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LOMS_CancelOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CancelOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LOMSServer).CancelOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loms.LOMS/CancelOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LOMSServer).CancelOrder(ctx, req.(*CancelOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LOMS_Stocks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StocksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LOMSServer).Stocks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loms.LOMS/Stocks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LOMSServer).Stocks(ctx, req.(*StocksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LOMS_ServiceDesc is the grpc.ServiceDesc for LOMS service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LOMS_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "loms.LOMS",
	HandlerType: (*LOMSServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateOrder",
			Handler:    _LOMS_CreateOrder_Handler,
		},
		{
			MethodName: "ListOrder",
			Handler:    _LOMS_ListOrder_Handler,
		},
		{
			MethodName: "OrderPaid",
			Handler:    _LOMS_OrderPaid_Handler,
		},
		{
			MethodName: "CancelOrder",
			Handler:    _LOMS_CancelOrder_Handler,
		},
		{
			MethodName: "Stocks",
			Handler:    _LOMS_Stocks_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "loms.proto",
}