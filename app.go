package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"memory-lane/app/papaya"
	"memory-lane/app/wingman"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

const PROTOCOL_ID = "p2p"

func main() {
	l := log.New(os.Stdout, "memory-lane ", log.LstdFlags)
	g, err := papaya.NewGallery(l)
	if err != nil {
		l.Fatal("error while instantiating gallery: ", err)
	}
	l.Println(g)

	// start a libp2p node
	node, err := libp2p.New()
	if err != nil {
		l.Fatalln(err)
	}
	// defer the close of the network connection
	defer node.Close()

	// Extract private address to send to Galactus
	multiAddr := newMultiAddr(node, l)
	node.SetStreamHandler(PROTOCOL_ID, func(s network.Stream) {
		handler := wingman.NewWingmanHandler(l)
		handler.HandleStream(s)
	})
	l.Println("Listening on:", multiAddr)

	peerAddr := flag.String("peer-address", "", "peer address")
	flag.Parse()
	if *peerAddr != "" {
		// Parse the multiaddr string.
		peerMA, err := multiaddr.NewMultiaddr(*peerAddr)
		if err != nil {
			l.Fatalf("failed parsing to peerMA: %v", err)
		}
		peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
		if err != nil {
			l.Fatalf("failed parsing to peer address info: %v", err)
		}

		// Connect to the node at the given address
		if err := node.Connect(context.Background(), *peerAddrInfo); err != nil {
			panic(err)
		}
		l.Println("Connected to:", peerAddrInfo.String())

		// Open a new stream to a connected node
		s, err := node.NewStream(context.Background(), peerAddrInfo.ID, PROTOCOL_ID)
		if err != nil {
			l.Fatalf("failed opening a new stream: %v", err)
		}

		// Encode JSON data and send over stream
		d := wingman.WingmanMessage{Message: "Hello, world!"}
		encoder := json.NewEncoder(s)
		if err := encoder.Encode(&d); err != nil {
			l.Fatalf("failed encoding message: %v", err)
		}
	}

	// wait for a SIGINT or SIGTERM signal (ctrl + c)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	// shut the node down
	if err := node.Close(); err != nil {
		panic(err)
	}
}

func newMultiAddr(node host.Host, l *log.Logger) string {
	var privateAddrs []multiaddr.Multiaddr

	for _, addr := range node.Addrs() {
		addrStr := addr.String()
		addrSplit := strings.Split(addrStr, "/")
		if addrSplit[1] == "ip4" && addrSplit[3] == "tcp" {
			ip := net.ParseIP(addrSplit[2])
			if ip != nil && ip.IsPrivate() {
				privateAddrs = append(privateAddrs, addr)
			}
		}
	}

	multiaddr := fmt.Sprintf("%s/%s/%s", privateAddrs[0].String(), PROTOCOL_ID, node.ID().String())

	return multiaddr
}
