package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"memory-lane/app/raccoon"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/libp2p/go-libp2p"
	network "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

type CrdtStruct struct {
	Name     string
	Commands []string
	Image    []byte
}

const protocolID = "pingPongCounter"

// TODO: How to connect a front end?
// TODO: contact galactus (http request) to register this node's multiaddr and Peer ID
// TODO: retrieve other nodes multiaddr + Peer ID from Galactus and connect to them
// TODO: perform necessary operations/routines with connected nodes (e.g. sync the albums, add/delete images etc.)
// TODO: probably need some better error handling
func main() {
	l := log.New(os.Stdout, "memory-lane ", log.LstdFlags)
	crdt := raccoon.NewCRDT(l)
	l.Println("CRDT initialized", crdt)

	// Add -peer-address flag (this is currently given as a command line argument but will be provided by Galactus)
	peerAddr := flag.String("peer-address", "", "peer address")
	flag.Parse()

	// start a libp2p node
	node, err := libp2p.New()
	if err != nil {
		panic(err)
	}
	// defer the close of the network connection
	defer node.Close()

	addrs := node.Addrs()
	fmt.Println("Listening on (IP-multiaddrs):")
	for _, addr := range addrs {
		fmt.Printf("\t%s\n", addr.String())
	}

	// Extract private address to send to Galactus
	var privateAddrs []multiaddr.Multiaddr
	for _, addr := range addrs {
		addrStr := addr.String()
		addrSplit := strings.Split(addrStr, "/")
		if addrSplit[1] == "ip4" && addrSplit[3] == "tcp" {
			ip := net.ParseIP(addrSplit[2])
			if ip != nil && ip.IsPrivate() {
				privateAddrs = append(privateAddrs, addr)
				fmt.Println("Private Address:", addrStr)
			}
		}
	}
	fmt.Println("Peer ID:", node.ID())
	fmt.Println("Adress to connect to:", privateAddrs[0].String()+"/p2p/"+node.ID().String())
	// Setup Stream Handlers
	// This gets called when peer connects and opens a stream to this node
	node.SetStreamHandler(protocolID, func(s network.Stream) {
		go readCRDT(s)
		crdtExample := CrdtStruct{
			Name:     fmt.Sprintf("CRDT of PeerId %s", node.ID()),
			Commands: []string{"delete photo"},
			Image:    []byte{0x01, 0x02, 0x03, 0x04},
		}
		go writeCRDT(s, crdtExample)
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
		go readCRDT(s)
		crdtExample := CrdtStruct{
			Name:     fmt.Sprintf("CRDT of PeerId %s", node.ID()),
			Commands: []string{"add photo", "delete photo", "delete photo"},
			Image:    []byte{0x01, 0x02, 0x03, 0x04},
		}
		go writeCRDT(s, crdtExample)
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

/* Util Methods */

// Send and Receive Data
func writeCRDT(s network.Stream, crdt CrdtStruct) {
	err := sendCRDT(crdt, s)
	if err != nil {
		panic(err)
	}
}

func readCRDT(s network.Stream) {
	fmt.Println("Reading from stream")
	for {
		crdt, err := receiveCRDT(s)
		if err != nil {
			panic(err)
		}
		// todo: process crdt
		fmt.Println("Processing CRDT", crdt)
	}
}

func sendCRDT(crdt CrdtStruct, s network.Stream) error {
	data, err := json.Marshal(crdt)
	if err != nil {
		return err
	}
	_, err = s.Write(data)
	fmt.Printf("Sending %s\n", string(data))
	if err != nil {
		return err
	}
	return nil
}

func receiveCRDT(s network.Stream) (CrdtStruct, error) {
	var crdt CrdtStruct
	err := json.NewDecoder(s).Decode(&crdt)
	if err != nil {
		return crdt, err
	}
	jsonCRDT, _ := json.Marshal(crdt)
	fmt.Printf("Received %s\n", jsonCRDT)
	return crdt, nil
}
