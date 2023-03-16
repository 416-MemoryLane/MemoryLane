package main

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/libp2p/go-libp2p"
	network "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

const protocolID = "pingPongCounter"

// TODO: How to connect a front end?
// TODO: contact galactus (http request) to register this node's multiaddr and Peer ID
// TODO: retrieve other nodes multiaddr + Peer ID from Galactus and connect to them
// TODO: perform necessary operations/routines with connected nodes (e.g. sync the albums, add/delete images etc.)
// TODO: probably need some better error handling
func main() {
	// Add -peer-address flag (this is currently given as a command line argument but will be provided by Galactus)
	peerAddr := flag.String("peer-address", "", "peer address")
	flag.Parse()

	// start a libp2p node with default settings and restricting it to IPv4 over 127.0.0.1 (only can communicate with nodes on same machine)
	// note:
	//	- that the port is set to 0, which means that the OS will assign a random port
	//  - the ip is set to 127.0.0.1 but will be the ip of the node in the network (public ip of computer)
	node, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"))
	if err != nil {
		panic(err)
	}
	// defer the close of the network connection
	defer node.Close()

	fmt.Println("Listening addresses (IP-multiaddr):", node.Addrs())
	fmt.Println("Peer ID:", node.ID())

	// Setup Stream Handlers
	// This gets called when peer connects and opens a stream to this node
	node.SetStreamHandler(protocolID, func(s network.Stream) {
		go writeCounter(s)
		go readCounter(s)
	})

	// Todo: remove the peer-address flag and use Galactus information for connection
	// Connect to peer if peer address is provided as command line argument
	if *peerAddr != "" {
		// Parse the multiaddr string.
		peerMA, err := multiaddr.NewMultiaddr(*peerAddr)
		if err != nil {
			panic(err)
		}
		peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
		if err != nil {
			panic(err)
		}

		// Connect to the node at the given address
		if err := node.Connect(context.Background(), *peerAddrInfo); err != nil {
			panic(err)
		}
		fmt.Println("Connected to", peerAddrInfo.String())

		// Open a stream with the given node
		s, err := node.NewStream(context.Background(), peerAddrInfo.ID, protocolID)
		if err != nil {
			panic(err)
		}

		// Start the write and read threads
		go writeCounter(s)
		go readCounter(s)
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

// Send and Receive Data
// Continue to send and receive a counter value until one of the nodes is killed
func writeCounter(s network.Stream) {
	var counter uint64

	for {
		<-time.After(time.Second)
		counter++

		err := binary.Write(s, binary.BigEndian, counter)
		if err != nil {
			panic(err)
		}
	}
}

func readCounter(s network.Stream) {
	for {
		var counter uint64

		err := binary.Read(s, binary.BigEndian, &counter)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Received %d from %s\n", counter, s.ID())
	}
}
