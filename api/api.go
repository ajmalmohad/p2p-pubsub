package api

import (
	"context"
	apigen "lumina/gen/api"
)

type ApiServer struct {
	apigen.UnimplementedApiServer
}

func (api *ApiServer) SendMessage(ctx context.Context, req *apigen.SendMessageRequest) (*apigen.SendMessageReply, error) {
	print(req.Value)
	return &apigen.SendMessageReply{Success: true}, nil
}

func (api *ApiServer) JoinRoom(ctx context.Context, req *apigen.JoinRoomRequest) (*apigen.JoinRoomReply, error) {
	return &apigen.JoinRoomReply{Success: true}, nil
}

func NewServer() *ApiServer {
	s := &ApiServer{}
	return s
}
