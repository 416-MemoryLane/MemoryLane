package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"memory-lane/app/galactus_client"
	"memory-lane/app/papaya"
	"memory-lane/app/wingman"
	"net"
	"net/http"
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
const GALLERY_DIR = "./memory-lane-gallery"

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
	g, err := papaya.NewGallery(GALLERY_DIR, l)
	if err != nil {
		l.Fatal("error while instantiating gallery: ", err)
	}
	l.Println("Gallery instantiated")

	// Get public IP address
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		l.Fatal("error while getting ip: ", err)
		return
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Fatal("error while reading ip:", err)
		return
	}

	// Start a libp2p node
	node, err := libp2p.New()
	if err != nil {
		l.Fatalln(err)
	}
	defer node.Close()

	// Extract multiaddr to send to Galactus
	maddr := newMultiAddr(string(ip), node, l)
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

	ticker := time.NewTicker(3 * time.Second)
	for range ticker.C {
		// Sync with Galactus
		syncResp, err := gc.Sync()
		if err != nil {
			l.Fatalf("Error syncing with Galactus: %v", err)
		}
		l.Printf("Synced successfully with Galactus for %v albums\n", len(*syncResp))

		// Reconcile gallery state
		albumIds, err := g.GetAlbumIDs()
		if err != nil {
			l.Fatalf("failed retrieving album IDs: %v", err)
		}

		peerAddrsToAlbums := map[string]*[]string{}
		for _, syncAlbum := range *syncResp {
			// Initialize any albums missing in local filesystem
			syncAlbumId := syncAlbum.AlbumID
			if !(*albumIds)[syncAlbumId] {
				_, err := g.AddAlbum(syncAlbumId, syncAlbum.AlbumName)
				if err != nil {
					l.Fatalf("failed reconciling gallery state: %v", err)
				}
			}

			// Add to map of peer addresses to albums
			for _, u := range syncAlbum.AuthorizedUsers {
				if maddr != u {
					_, ok := peerAddrsToAlbums[u]
					if !ok {
						peerAddrsToAlbums[u] = &[]string{}
					}
					as := peerAddrsToAlbums[u]

					newAlbums := append(*as, syncAlbumId)
					peerAddrsToAlbums[u] = &newAlbums
				}

			}
		}

		for addr, albums := range peerAddrsToAlbums {
			// Parse the multiaddr string.
			peerMA, err := multiaddr.NewMultiaddr(addr)
			if err != nil {
				l.Printf("failed parsing to peerMA: %v", err)
				continue
			}
			peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
			if err != nil {
				l.Printf("failed parsing to peer address info: %v", err)
				continue
			}

			// Connect to the node at the given address
			if err := node.Connect(context.Background(), *peerAddrInfo); err != nil {
				l.Printf("failed to connect to peer: %v", err)
				continue
			}
			l.Println("Connected to:", peerAddrInfo.String())

			// Open a new stream to a connected node
			s, err := node.NewStream(context.Background(), peerAddrInfo.ID, PROTOCOL_ID)
			if err != nil {
				l.Fatalf("failed opening a new stream: %v", err)
			}

			// Send message to each album
			for _, aid := range *albums {
				crdt, err := g.GetAlbumCRDT(aid)
				if err != nil {
					l.Fatalf("failed retrieving album CRDT: %v", err)
				}

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
					} else {
						l.Printf("sent msg to: %v\n for album: %v\n", peerAddrInfo.String(), aid)
					}
				}()
			}
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

func newMultiAddr(ipStr string, node host.Host, l *log.Logger) string {
	var privateAddrs []multiaddr.Multiaddr

	for _, addr := range node.Addrs() {
		addrStr := addr.String()
		addrSplit := strings.Split(addrStr, "/")
		if addrSplit[1] == "ip4" && addrSplit[3] == "tcp" {
			ip := net.ParseIP(addrSplit[2])
			addrSplit[2] = ipStr
			multiaddrStr := strings.Join(addrSplit, "/")
			mAddr, err := multiaddr.NewMultiaddr(multiaddrStr)
			if err != nil {
				l.Fatalf("error parsing multiaddr: %v", err)
			}
			if ip != nil && ip.IsPrivate() {
				privateAddrs = append(privateAddrs, mAddr)
			}
		}
	}

	multiaddr := fmt.Sprintf("%s/%s/%s", privateAddrs[0].String(), PROTOCOL_ID, node.ID().String())

	return multiaddr
}
