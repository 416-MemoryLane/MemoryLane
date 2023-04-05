package wingman

import (
	"log"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

func MultiaddrStrToPeerAddrInfo(mid string, l *log.Logger) *peer.AddrInfo {
	// Parse the multiaddr string.
	peerMA, err := multiaddr.NewMultiaddr(mid)
	if err != nil {
		l.Fatalf("failed parsing to peerMA: %v", err)
	}
	peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
	if err != nil {
		l.Fatalf("failed parsing to peer address info: %v", err)
	}

	return peerAddrInfo
}
