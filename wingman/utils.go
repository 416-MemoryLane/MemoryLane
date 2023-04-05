package wingman

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// Convert a multiaddr string to *peer.AddrInfo
func (wh *WingmanHandler) MultiaddrStrToPeerAddrInfo(mid string) *peer.AddrInfo {
	// Parse the multiaddr string.
	peerMA, err := multiaddr.NewMultiaddr(mid)
	if err != nil {
		wh.l.Fatalf("failed parsing to peerMA: %v", err)
	}
	peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
	if err != nil {
		wh.l.Fatalf("failed parsing to peer address info: %v", err)
	}

	return peerAddrInfo
}
