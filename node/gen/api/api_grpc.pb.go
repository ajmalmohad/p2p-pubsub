// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

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

// ApiClient is the client API for Api service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ApiClient interface {
	SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageReply, error)
	GetRoomParticipants(ctx context.Context, in *GetRoomParticipantsRequest, opts ...grpc.CallOption) (*GetRoomParticipantsResponse, error)
	SubscribeEvents(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (Api_SubscribeEventsClient, error)
}

type apiClient struct {
	cc grpc.ClientConnInterface
}

func NewApiClient(cc grpc.ClientConnInterface) ApiClient {
	return &apiClient{cc}
}

func (c *apiClient) SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageReply, error) {
	out := new(SendMessageReply)
	err := c.cc.Invoke(ctx, "/api.Api/SendMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiClient) GetRoomParticipants(ctx context.Context, in *GetRoomParticipantsRequest, opts ...grpc.CallOption) (*GetRoomParticipantsResponse, error) {
	out := new(GetRoomParticipantsResponse)
	err := c.cc.Invoke(ctx, "/api.Api/GetRoomParticipants", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiClient) SubscribeEvents(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (Api_SubscribeEventsClient, error) {
	stream, err := c.cc.NewStream(ctx, &Api_ServiceDesc.Streams[0], "/api.Api/SubscribeEvents", opts...)
	if err != nil {
		return nil, err
	}
	x := &apiSubscribeEventsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Api_SubscribeEventsClient interface {
	Recv() (*Event, error)
	grpc.ClientStream
}

type apiSubscribeEventsClient struct {
	grpc.ClientStream
}

func (x *apiSubscribeEventsClient) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ApiServer is the server API for Api service.
// All implementations must embed UnimplementedApiServer
// for forward compatibility
type ApiServer interface {
	SendMessage(context.Context, *SendMessageRequest) (*SendMessageReply, error)
	GetRoomParticipants(context.Context, *GetRoomParticipantsRequest) (*GetRoomParticipantsResponse, error)
	SubscribeEvents(*SubscribeRequest, Api_SubscribeEventsServer) error
	mustEmbedUnimplementedApiServer()
}

// UnimplementedApiServer must be embedded to have forward compatible implementations.
type UnimplementedApiServer struct {
}

func (UnimplementedApiServer) SendMessage(context.Context, *SendMessageRequest) (*SendMessageReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedApiServer) GetRoomParticipants(context.Context, *GetRoomParticipantsRequest) (*GetRoomParticipantsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoomParticipants not implemented")
}
func (UnimplementedApiServer) SubscribeEvents(*SubscribeRequest, Api_SubscribeEventsServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeEvents not implemented")
}
func (UnimplementedApiServer) mustEmbedUnimplementedApiServer() {}

// UnsafeApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ApiServer will
// result in compilation errors.
type UnsafeApiServer interface {
	mustEmbedUnimplementedApiServer()
}

func RegisterApiServer(s grpc.ServiceRegistrar, srv ApiServer) {
	s.RegisterService(&Api_ServiceDesc, srv)
}

func _Api_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Api/SendMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServer).SendMessage(ctx, req.(*SendMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Api_GetRoomParticipants_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoomParticipantsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServer).GetRoomParticipants(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Api/GetRoomParticipants",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServer).GetRoomParticipants(ctx, req.(*GetRoomParticipantsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Api_SubscribeEvents_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ApiServer).SubscribeEvents(m, &apiSubscribeEventsServer{stream})
}

type Api_SubscribeEventsServer interface {
	Send(*Event) error
	grpc.ServerStream
}

type apiSubscribeEventsServer struct {
	grpc.ServerStream
}

func (x *apiSubscribeEventsServer) Send(m *Event) error {
	return x.ServerStream.SendMsg(m)
}

// Api_ServiceDesc is the grpc.ServiceDesc for Api service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Api_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.Api",
	HandlerType: (*ApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMessage",
			Handler:    _Api_SendMessage_Handler,
		},
		{
			MethodName: "GetRoomParticipants",
			Handler:    _Api_GetRoomParticipants_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeEvents",
			Handler:       _Api_SubscribeEvents_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/api.proto",
}
