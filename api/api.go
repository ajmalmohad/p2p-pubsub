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
	print("Message Sent\n") //
	return &apigen.SendMessageReply{Success: true}, nil
}

func (api *ApiServer) GetRoomParticipants(ctx context.Context, req *apigen.GetRoomParticipantsRequest) (*apigen.GetRoomParticipantsResponse, error) {
	var participants []*apigen.RoomParticipant
	allparts := api.cr.ListPeers()
	print("Room Participants: \n") //
	for _, part := range allparts {
		participants = append(participants, &apigen.RoomParticipant{
			Id:       part.String(),
			Nickname: "I don't Know",
		})
		print(participants[len(participants)-1].Id, ": ", participants[len(participants)-1].Nickname, "\n") //
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

			PrintMessage(&message) //

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

			print("Peer Joined: \n") //
			print(m.PeerID, "\n")    //

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

			print("Peer Left: \n") //
			print(m.PeerID, "\n")  //

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

func PrintMessage(m *apigen.ChatMessage) {
	print("Message: \n")
	print("Sender ID: ", m.SenderId, "\n")
	print("NickName: ", m.SenderNickname, "\n")
	print("Timestamp: ", m.Timestamp, "\n")
	print("Value: ", m.Value, "\n")
}
