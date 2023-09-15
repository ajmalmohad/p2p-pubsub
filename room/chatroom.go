package room

import (
	"context"
	"encoding/json"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type ChatRoom struct {
	Messages chan *ChatMessage
	PeerJoin chan *PeerJoin
	PeerLeft chan *PeerLeft

	ctx   context.Context
	ps    *pubsub.PubSub
	topic *pubsub.Topic
	sub   *pubsub.Subscription

	roomName string
	self     peer.ID
	nick     string
}

type ChatMessage struct {
	Message    string
	SenderID   string
	SenderNick string
	Timestamp  string
}

type PeerJoin struct {
	PeerID string
}

type PeerLeft struct {
	PeerID string
}

func JoinChatRoom(ctx context.Context, ps *pubsub.PubSub, selfID peer.ID, nickname string, roomName string) (*ChatRoom, error) {
	topic, err := ps.Join(topicName(roomName))
	if err != nil {
		return nil, err
	}

	sub, err := topic.Subscribe()
	if err != nil {
		return nil, err
	}

	cr := &ChatRoom{
		ctx:      ctx,
		ps:       ps,
		topic:    topic,
		sub:      sub,
		self:     selfID,
		nick:     nickname,
		roomName: roomName,
		Messages: make(chan *ChatMessage, 128),
		PeerJoin: make(chan *PeerJoin, 128),
		PeerLeft: make(chan *PeerLeft, 128),
	}

	go cr.readLoop()
	go handlePeerEvents(cr)
	return cr, nil
}

func handlePeerEvents(cr *ChatRoom) {
	handler, err := cr.topic.EventHandler()
	if err != nil {
		panic(err)
	}

	for {
		peerEvt, err := handler.NextPeerEvent(context.Background())
		if err != nil {
			print("failed receiving room topic peer event")
			continue
		}

		switch peerEvt.Type {
		case pubsub.PeerLeave:
			cr.PeerLeft <- &PeerLeft{
				PeerID: peerEvt.Peer.Pretty(),
			}
		case pubsub.PeerJoin:
			cr.PeerJoin <- &PeerJoin{
				PeerID: peerEvt.Peer.Pretty(),
			}
		}
	}
}

func (cr *ChatRoom) Publish(message string) error {
	m := ChatMessage{
		Message:    message,
		SenderID:   cr.self.Pretty(),
		SenderNick: cr.nick,
		Timestamp:  time.Now().String(),
	}
	msgBytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return cr.topic.Publish(cr.ctx, msgBytes)
}

func (cr *ChatRoom) ListPeers() []peer.ID {
	return cr.ps.ListPeers(topicName(cr.roomName))
}

func (cr *ChatRoom) readLoop() {
	for {
		msg, err := cr.sub.Next(cr.ctx)
		if err != nil {
			close(cr.Messages)
			return
		}

		if msg.ReceivedFrom == cr.self {
			continue
		}

		cm := new(ChatMessage)
		err = json.Unmarshal(msg.Data, cm)
		if err != nil {
			continue
		}
		cr.Messages <- cm
	}
}

func topicName(roomName string) string {
	return "chat-room:" + roomName
}
