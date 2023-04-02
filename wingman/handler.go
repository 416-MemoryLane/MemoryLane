package wingman

import "log"

type WingmanHandler struct {
	l *log.Logger
}

func NewWingmanHandler(l *log.Logger) *WingmanHandler {
	return &WingmanHandler{l}
}
