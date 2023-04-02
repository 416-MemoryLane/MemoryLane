package wingman

import (
	"encoding/json"

	"github.com/libp2p/go-libp2p/core/network"
)

func (wh *WingmanHandler) HandleStream(stream network.Stream) {
	defer stream.Close()

	decoder := json.NewDecoder(stream)
	var d WingmanMessage
	if err := decoder.Decode(&d); err != nil {
		wh.l.Printf("error handling stream: %v", err)
		return
	}

	wh.l.Println("received message:", d.Message)
}
