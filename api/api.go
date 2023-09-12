package api

import (
	"context"
	apigen "lumina/gen/api"
	chatroom "lumina/room"
)

type ApiServer struct {
	apigen.UnimplementedApiServer
	cr *chatroom.ChatRoom
}

func (api *ApiServer) SendMessage(ctx context.Context, req *apigen.SendMessageRequest) (*apigen.SendMessageReply, error) {
	api.cr.Publish(req.Value)
	return &apigen.SendMessageReply{Success: true}, nil
}

func (api *ApiServer) GetRoomParticipants(ctx context.Context, req *apigen.GetRoomParticipantsRequest) (*apigen.GetRoomParticipantsResponse, error) {
	var participants []*apigen.RoomParticipant
	allparts := api.cr.ListPeers()
	for _, part := range allparts {
		participants = append(participants, &apigen.RoomParticipant{
			Id:       part.String(),
			Nickname: "I don't Know",
		})
	}

	return &apigen.GetRoomParticipantsResponse{Participants: participants}, nil
}

func NewServer(cr *chatroom.ChatRoom) *ApiServer {
	s := &ApiServer{cr: cr}
	return s
}
