package api

import (
	"context"
	apigen "lumina/gen/api"
)

type ApiServer struct {
	apigen.UnimplementedApiServer
}

func (api *ApiServer) SendMessage(context.Context, *apigen.MessageRequest) (*apigen.MessageReply, error) {
	return &apigen.MessageReply{Sent: true}, nil
}

func NewServer() *ApiServer {
	s := &ApiServer{}
	return s
}
