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

func (api *ApiServer) GetRoomParticipants(ctx context.Context, req *apigen.GetRoomParticipantsRequest) (*apigen.GetRoomParticipantsResponse, error) {
	var participants []*apigen.RoomParticipant
	participants = append(participants, &apigen.RoomParticipant{
		Id:       "1",
		Nickname: "Ajmal",
	})

	participants = append(participants, &apigen.RoomParticipant{
		Id:       "2",
		Nickname: "Lamja",
	})

	return &apigen.GetRoomParticipantsResponse{Participants: participants}, nil
}

func NewServer() *ApiServer {
	s := &ApiServer{}
	return s
}
