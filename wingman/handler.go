package wingman

import (
	"log"
	"memory-lane/app/papaya"

	"github.com/libp2p/go-libp2p/core/host"
)

type WingmanHandler struct {
	Multiaddr  string
	ProtocolId string
	Node       *host.Host
	Gallery    *papaya.Gallery
	l          *log.Logger
}

func NewWingmanHandler(
	multiAddr string,
	pid string,
	n *host.Host,
	g *papaya.Gallery,
	l *log.Logger,
) *WingmanHandler {
	return &WingmanHandler{multiAddr, pid, n, g, l}
}
