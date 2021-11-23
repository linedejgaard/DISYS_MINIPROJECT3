// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package Auction

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

// AuctionServiceClient is the client API for AuctionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuctionServiceClient interface {
	Election(ctx context.Context, in *ElectionRequest, opts ...grpc.CallOption) (*ElectionReply, error)
	LeaderDeclaration(ctx context.Context, in *LeaderDeclarationRequest, opts ...grpc.CallOption) (*LeaderDeclarationReply, error)
	Join(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (*JoinReply, error)
	Leave(ctx context.Context, in *LeaveRequest, opts ...grpc.CallOption) (*LeaveReply, error)
	UpdatePorts(ctx context.Context, in *UpdatePortsRequest, opts ...grpc.CallOption) (*UpdatePortsReply, error)
	UpdateActionStatus(ctx context.Context, in *UpdateAuctionStatusRequest, opts ...grpc.CallOption) (*UpdateActionStatusReply, error)
	Bid(ctx context.Context, in *BidRequest, opts ...grpc.CallOption) (*BidReply, error)
	PublishResult(ctx context.Context, in *PublishResultRequest, opts ...grpc.CallOption) (*PublishResultReply, error)
	GetState(ctx context.Context, in *GetStateRequest, opts ...grpc.CallOption) (*GetStateReply, error)
	MakeNewAuction(ctx context.Context, in *MakeNewAuctionRequest, opts ...grpc.CallOption) (*MakeNewAuctionReply, error)
}

type auctionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuctionServiceClient(cc grpc.ClientConnInterface) AuctionServiceClient {
	return &auctionServiceClient{cc}
}

func (c *auctionServiceClient) Election(ctx context.Context, in *ElectionRequest, opts ...grpc.CallOption) (*ElectionReply, error) {
	out := new(ElectionReply)
	err := c.cc.Invoke(ctx, "/Auction.AuctionService/Election", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auctionServiceClient) LeaderDeclaration(ctx context.Context, in *LeaderDeclarationRequest, opts ...grpc.CallOption) (*LeaderDeclarationReply, error) {
	out := new(LeaderDeclarationReply)
	err := c.cc.Invoke(ctx, "/Auction.AuctionService/LeaderDeclaration", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auctionServiceClient) Join(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (*JoinReply, error) {
	out := new(JoinReply)
	err := c.cc.Invoke(ctx, "/Auction.AuctionService/Join", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auctionServiceClient) Leave(ctx context.Context, in *LeaveRequest, opts ...grpc.CallOption) (*LeaveReply, error) {
	out := new(LeaveReply)
	err := c.cc.Invoke(ctx, "/Auction.AuctionService/Leave", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auctionServiceClient) UpdatePorts(ctx context.Context, in *UpdatePortsRequest, opts ...grpc.CallOption) (*UpdatePortsReply, error) {
	out := new(UpdatePortsReply)
	err := c.cc.Invoke(ctx, "/Auction.AuctionService/UpdatePorts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auctionServiceClient) UpdateActionStatus(ctx context.Context, in *UpdateAuctionStatusRequest, opts ...grpc.CallOption) (*UpdateActionStatusReply, error) {
	out := new(UpdateActionStatusReply)
	err := c.cc.Invoke(ctx, "/Auction.AuctionService/UpdateActionStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auctionServiceClient) Bid(ctx context.Context, in *BidRequest, opts ...grpc.CallOption) (*BidReply, error) {
	out := new(BidReply)
	err := c.cc.Invoke(ctx, "/Auction.AuctionService/Bid", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auctionServiceClient) PublishResult(ctx context.Context, in *PublishResultRequest, opts ...grpc.CallOption) (*PublishResultReply, error) {
	out := new(PublishResultReply)
	err := c.cc.Invoke(ctx, "/Auction.AuctionService/PublishResult", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auctionServiceClient) GetState(ctx context.Context, in *GetStateRequest, opts ...grpc.CallOption) (*GetStateReply, error) {
	out := new(GetStateReply)
	err := c.cc.Invoke(ctx, "/Auction.AuctionService/GetState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auctionServiceClient) MakeNewAuction(ctx context.Context, in *MakeNewAuctionRequest, opts ...grpc.CallOption) (*MakeNewAuctionReply, error) {
	out := new(MakeNewAuctionReply)
	err := c.cc.Invoke(ctx, "/Auction.AuctionService/MakeNewAuction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuctionServiceServer is the server API for AuctionService service.
// All implementations must embed UnimplementedAuctionServiceServer
// for forward compatibility
type AuctionServiceServer interface {
	Election(context.Context, *ElectionRequest) (*ElectionReply, error)
	LeaderDeclaration(context.Context, *LeaderDeclarationRequest) (*LeaderDeclarationReply, error)
	Join(context.Context, *JoinRequest) (*JoinReply, error)
	Leave(context.Context, *LeaveRequest) (*LeaveReply, error)
	UpdatePorts(context.Context, *UpdatePortsRequest) (*UpdatePortsReply, error)
	UpdateActionStatus(context.Context, *UpdateAuctionStatusRequest) (*UpdateActionStatusReply, error)
	Bid(context.Context, *BidRequest) (*BidReply, error)
	PublishResult(context.Context, *PublishResultRequest) (*PublishResultReply, error)
	GetState(context.Context, *GetStateRequest) (*GetStateReply, error)
	MakeNewAuction(context.Context, *MakeNewAuctionRequest) (*MakeNewAuctionReply, error)
	mustEmbedUnimplementedAuctionServiceServer()
}

// UnimplementedAuctionServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAuctionServiceServer struct {
}

func (UnimplementedAuctionServiceServer) Election(context.Context, *ElectionRequest) (*ElectionReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Election not implemented")
}
func (UnimplementedAuctionServiceServer) LeaderDeclaration(context.Context, *LeaderDeclarationRequest) (*LeaderDeclarationReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LeaderDeclaration not implemented")
}
func (UnimplementedAuctionServiceServer) Join(context.Context, *JoinRequest) (*JoinReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Join not implemented")
}
func (UnimplementedAuctionServiceServer) Leave(context.Context, *LeaveRequest) (*LeaveReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Leave not implemented")
}
func (UnimplementedAuctionServiceServer) UpdatePorts(context.Context, *UpdatePortsRequest) (*UpdatePortsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePorts not implemented")
}
func (UnimplementedAuctionServiceServer) UpdateActionStatus(context.Context, *UpdateAuctionStatusRequest) (*UpdateActionStatusReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateActionStatus not implemented")
}
func (UnimplementedAuctionServiceServer) Bid(context.Context, *BidRequest) (*BidReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Bid not implemented")
}
func (UnimplementedAuctionServiceServer) PublishResult(context.Context, *PublishResultRequest) (*PublishResultReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishResult not implemented")
}
func (UnimplementedAuctionServiceServer) GetState(context.Context, *GetStateRequest) (*GetStateReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetState not implemented")
}
func (UnimplementedAuctionServiceServer) MakeNewAuction(context.Context, *MakeNewAuctionRequest) (*MakeNewAuctionReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MakeNewAuction not implemented")
}
func (UnimplementedAuctionServiceServer) mustEmbedUnimplementedAuctionServiceServer() {}

// UnsafeAuctionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuctionServiceServer will
// result in compilation errors.
type UnsafeAuctionServiceServer interface {
	mustEmbedUnimplementedAuctionServiceServer()
}

func RegisterAuctionServiceServer(s grpc.ServiceRegistrar, srv AuctionServiceServer) {
	s.RegisterService(&AuctionService_ServiceDesc, srv)
}

func _AuctionService_Election_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ElectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServiceServer).Election(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Auction.AuctionService/Election",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServiceServer).Election(ctx, req.(*ElectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuctionService_LeaderDeclaration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaderDeclarationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServiceServer).LeaderDeclaration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Auction.AuctionService/LeaderDeclaration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServiceServer).LeaderDeclaration(ctx, req.(*LeaderDeclarationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuctionService_Join_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServiceServer).Join(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Auction.AuctionService/Join",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServiceServer).Join(ctx, req.(*JoinRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuctionService_Leave_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServiceServer).Leave(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Auction.AuctionService/Leave",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServiceServer).Leave(ctx, req.(*LeaveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuctionService_UpdatePorts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePortsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServiceServer).UpdatePorts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Auction.AuctionService/UpdatePorts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServiceServer).UpdatePorts(ctx, req.(*UpdatePortsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuctionService_UpdateActionStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAuctionStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServiceServer).UpdateActionStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Auction.AuctionService/UpdateActionStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServiceServer).UpdateActionStatus(ctx, req.(*UpdateAuctionStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuctionService_Bid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BidRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServiceServer).Bid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Auction.AuctionService/Bid",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServiceServer).Bid(ctx, req.(*BidRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuctionService_PublishResult_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishResultRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServiceServer).PublishResult(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Auction.AuctionService/PublishResult",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServiceServer).PublishResult(ctx, req.(*PublishResultRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuctionService_GetState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServiceServer).GetState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Auction.AuctionService/GetState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServiceServer).GetState(ctx, req.(*GetStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuctionService_MakeNewAuction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MakeNewAuctionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServiceServer).MakeNewAuction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Auction.AuctionService/MakeNewAuction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServiceServer).MakeNewAuction(ctx, req.(*MakeNewAuctionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuctionService_ServiceDesc is the grpc.ServiceDesc for AuctionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuctionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Auction.AuctionService",
	HandlerType: (*AuctionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Election",
			Handler:    _AuctionService_Election_Handler,
		},
		{
			MethodName: "LeaderDeclaration",
			Handler:    _AuctionService_LeaderDeclaration_Handler,
		},
		{
			MethodName: "Join",
			Handler:    _AuctionService_Join_Handler,
		},
		{
			MethodName: "Leave",
			Handler:    _AuctionService_Leave_Handler,
		},
		{
			MethodName: "UpdatePorts",
			Handler:    _AuctionService_UpdatePorts_Handler,
		},
		{
			MethodName: "UpdateActionStatus",
			Handler:    _AuctionService_UpdateActionStatus_Handler,
		},
		{
			MethodName: "Bid",
			Handler:    _AuctionService_Bid_Handler,
		},
		{
			MethodName: "PublishResult",
			Handler:    _AuctionService_PublishResult_Handler,
		},
		{
			MethodName: "GetState",
			Handler:    _AuctionService_GetState_Handler,
		},
		{
			MethodName: "MakeNewAuction",
			Handler:    _AuctionService_MakeNewAuction_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Auction/Auction.proto",
}