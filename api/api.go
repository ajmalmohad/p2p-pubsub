package api

import (
	"context"
	"log"
	apigen "lumina/gen/api"
	chatroom "lumina/room"
)

type ApiServer struct {
	apigen.UnimplementedApiServer
	cr *chatroom.ChatRoom
}

func (api *ApiServer) SendMessage(ctx context.Context, req *apigen.SendMessageRequest) (*apigen.SendMessageReply, error) {
	err := api.cr.Publish(req.Value)
	if err != nil {
		panic(err)
	}
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

func (api *ApiServer) SubscribeEvents(req *apigen.SubscribeRequest, stream apigen.Api_SubscribeEventsServer) error {
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case m := <-api.cr.Messages:
			message := apigen.ChatMessage{
				SenderId:       m.SenderID,
				SenderNickname: m.SenderNick,
				Timestamp:      m.Timestamp,
				Value:          m.Message,
			}

			event := apigen.Event{
				Type:    1,
				Message: &message,
			}

			err := stream.Send(&event)
			if err != nil {
				log.Println(err.Error())
			}
		case m := <-api.cr.PeerJoin:
			peer := apigen.PeerJoin{
				PeerId: m.PeerID,
			}

			event := apigen.Event{
				Type:     2,
				PeerJoin: &peer,
			}

			err := stream.Send(&event)
			if err != nil {
				log.Println(err.Error())
			}
		case m := <-api.cr.PeerLeft:
			peer := apigen.PeerLeft{
				PeerId: m.PeerID,
			}

			event := apigen.Event{
				Type:     3,
				PeerLeft: &peer,
			}

			err := stream.Send(&event)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
}

func NewServer(cr *chatroom.ChatRoom) *ApiServer {
	s := &ApiServer{cr: cr}
	return s
}
