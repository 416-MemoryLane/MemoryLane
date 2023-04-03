package wingman

import (
	"log"
	"memory-lane/app/papaya"
)

type WingmanHandler struct {
	Multiaddr string
	Gallery   *papaya.Gallery
	l         *log.Logger
}

func NewWingmanHandler(m string, g *papaya.Gallery, l *log.Logger) *WingmanHandler {
	return &WingmanHandler{m, g, l}
}
