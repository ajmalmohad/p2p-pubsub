package room

import (
	"context"
	"encoding/json"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// ChatRoomBufSize is the number of incoming messages to buffer for each topic.
const ChatRoomBufSize = 128
const PeerEventBufSize = 128

// ChatRoom represents a subscription to a single PubSub topic. Messages
// can be published to the topic with ChatRoom.Publish, and received
// messages are pushed to the Messages channel.
type ChatRoom struct {
	// Messages is a channel of messages received from other peers in the chat room
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

// ChatMessage gets converted to/from JSON and sent in the body of pubsub messages.
type ChatMessage struct {
	Message    string
	SenderID   string
	SenderNick string
	Timestamp  string
}

// Peer Join.
type PeerJoin struct {
	PeerID string
}

// Peer Left.
type PeerLeft struct {
	PeerID string
}

// JoinChatRoom tries to subscribe to the PubSub topic for the room name, returning
// a ChatRoom on success.
func JoinChatRoom(ctx context.Context, ps *pubsub.PubSub, selfID peer.ID, nickname string, roomName string) (*ChatRoom, error) {
	// join the pubsub topic
	topic, err := ps.Join(topicName(roomName))
	if err != nil {
		return nil, err
	}

	// and subscribe to it
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
		Messages: make(chan *ChatMessage, ChatRoomBufSize),
		PeerJoin: make(chan *PeerJoin, PeerEventBufSize),
		PeerLeft: make(chan *PeerLeft, PeerEventBufSize),
	}

	// start reading messages from the subscription in a loop
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

// Publish sends a message to the pubsub topic.
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

// readLoop pulls messages from the pubsub topic and pushes them onto the Messages channel.
func (cr *ChatRoom) readLoop() {
	for {
		msg, err := cr.sub.Next(cr.ctx)
		if err != nil {
			close(cr.Messages)
			return
		}
		// only forward messages delivered by others
		if msg.ReceivedFrom == cr.self {
			continue
		}

		print(msg.Data)

		cm := new(ChatMessage)
		err = json.Unmarshal(msg.Data, cm)
		if err != nil {
			continue
		}
		// send valid messages onto the Messages channel
		cr.Messages <- cm
	}
}

func topicName(roomName string) string {
	return "chat-room:" + roomName
}
