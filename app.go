package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"memory-lane/app/galactus_client"
	"memory-lane/app/papaya"
	"memory-lane/app/wingman"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

const PROTOCOL_ID = "p2p"
const GALACTUS_API = "https://memory-lane-381119.wl.r.appspot.com"

func main() {
	l := log.New(os.Stdout, "", log.Lshortfile|log.Ltime)

	// Define flags
	unFlag := "username"
	pwFlag := "password"
	unPtr := flag.String(unFlag, "", "Username")
	pwPtr := flag.String(pwFlag, "", "Password")

	// Parse command line arguments
	flag.Parse()

	// Access flag values
	un := *unPtr
	pw := *pwPtr

	// Check if required flags are provided
	if un == "" || pw == "" {
		l.Printf("Please provide both --%s and --%s flags\n", unFlag, pwFlag)
		return
	}

	// Instantiate Gallery
	g, err := papaya.NewGallery(l)
	if err != nil {
		l.Fatal("error while instantiating gallery: ", err)
	}
	l.Println("Gallery instantiated")

	// Start a libp2p node
	node, err := libp2p.New()
	if err != nil {
		l.Fatalln(err)
	}
	defer node.Close()

	// Extract multiaddr to send to Galactus
	maddr := newMultiAddr(node, l)
	node.SetStreamHandler(PROTOCOL_ID, func(s network.Stream) {
		handler := wingman.NewWingmanHandler(maddr, PROTOCOL_ID, &node, g, l)
		handler.HandleStream(s)
	})
	l.Println("Listening on:", maddr)

	// Instantiate Galactus Client and log in
	gc := galactus_client.NewGalactusClient(GALACTUS_API, un, pw, maddr, l)
	loginResp, err := gc.Login()
	if err != nil {
		l.Fatalf("Error logging in for user %s: %v", un, err)
	}
	l.Printf("%s", loginResp.Message)
	gc.AuthToken = loginResp.Token

	syncResp, err := gc.Sync()
	if err != nil {
		l.Fatalf("Error syncing with Galactus: %v", err)
	}
	l.Printf("Synced successfully with Galactus for %v albums\n", len(*syncResp))

	// TODO: should replace with multiaddrs received from Galactus
	peerAddrs := []string{
		// "/ip4/172.28.67.129/tcp/36891/p2p/12D3KooWMW1y5JcJ95DYJ7pssShb7F3bCWwWzvMZ81Yxa9jfQBbh",
		// "/ip4/172.28.67.129/tcp/34263/p2p/12D3KooWJZMXwYo8GpAbDiezXFPwViKnjWT7NKDwe3njnM8frn4t",
	}

	for _, addr := range peerAddrs {
		// Parse the multiaddr string.
		peerMA, err := multiaddr.NewMultiaddr(addr)
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

		// Retrieve album directories from filesystem and create a stream for each album
		albumCRDTs, err := g.GetAlbumCRDTs()
		if err != nil {
			l.Fatalf("failed retrieving album CRDTs: %v", err)
		}

		for _, crdt := range *albumCRDTs {
			aid := crdt.Album
			l.Println("Creating a stream for album:", aid)

			// Construct initial wingmanMsg
			wingmanMsg := wingman.WingmanMessage{
				SenderMultiAddr: maddr,
				Album:           aid,
				Crdt:            crdt,
				Photos:          nil,
			}

			go func() {
				// Encode JSON data and send over stream
				encoder := json.NewEncoder(s)
				if err := encoder.Encode(&wingmanMsg); err != nil {
					l.Fatalf("failed encoding message: %v", err)
				}

				ticker := time.NewTicker(3 * time.Second)
				for range ticker.C {
					crdt, err := g.GetAlbumCRDT(aid)
					if err != nil {
						l.Fatalf("failed retrieving crdt: %v", err)
					}

					wingmanMsg = wingman.WingmanMessage{
						SenderMultiAddr: maddr,
						Album:           aid,
						Crdt:            crdt,
						Photos:          nil,
					}

					if err = encoder.Encode(&wingmanMsg); err != nil {
						l.Printf("failed encoding message: %v\n", err)
						continue
					}

					l.Printf("sent msg to: %v\n for album: %v\n", peerAddrInfo.String(), aid)
				}
			}()
		}
	}

	// Gracefully shutdown node
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)

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
