package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"google.golang.org/grpc"

	api "lumina/api"
	apigen "lumina/gen/api"
	chatroom "lumina/room"
)

const DiscoveryInterval = time.Hour
const DiscoveryServiceTag = "lumina-pubsub"

func main() {
	nickFlag := flag.String("nick", "", "nickname to use in chat. will be generated if empty")
	roomFlag := flag.String("room", "lobby", "name of chat room to join")
	portFlag := flag.String("port", "3000", "port to open grpc server")
	flag.Parse()

	ctx := context.Background()

	node, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	if err != nil {
		panic(err)
	}

	ps, err := pubsub.NewGossipSub(ctx, node)
	if err != nil {
		panic(err)
	}

	nick := *nickFlag
	if len(nick) == 0 {
		nick = fmt.Sprintf("%s-%s", os.Getenv("USER"), shortID(node.ID()))
	}

	room := *roomFlag

	cr, err := chatroom.JoinChatRoom(ctx, ps, node.ID(), nick, room)
	if err != nil {
		panic(err)
	}

	if err := setupDiscovery(node); err != nil {
		panic(err)
	}

	port := *portFlag
	if len(port) == 0 {
		panic("No Port Specified")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		panic(err)
	}

	print("API started in localhost:", port, "\n")
	grpcServer := grpc.NewServer()
	apigen.RegisterApiServer(grpcServer, api.NewServer(cr))
	grpcServer.Serve(lis)

}

func shortID(p peer.ID) string {
	pretty := p.Pretty()
	return pretty[len(pretty)-8:]
}

type discoveryNotifee struct {
	h host.Host
}

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	print("Found a Peer on the network\n") //
	err := n.h.Connect(context.Background(), pi)
	if err != nil {
		fmt.Printf("error connecting to peer %s: %s\n", pi.ID.Pretty(), err)
	}
}

func setupDiscovery(h host.Host) error {
	s := mdns.NewMdnsService(h, DiscoveryServiceTag, &discoveryNotifee{h: h})
	return s.Start()
}
