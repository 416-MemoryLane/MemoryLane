package wingman

import (
	"log"
	"memory-lane/app/papaya"
)

type WingmanHandler struct {
	Gallery *papaya.Gallery
	l       *log.Logger
}

func NewWingmanHandler(g *papaya.Gallery, l *log.Logger) *WingmanHandler {
	return &WingmanHandler{g, l}
}
