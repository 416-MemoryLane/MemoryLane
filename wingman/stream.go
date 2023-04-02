package wingman

import (
	"encoding/json"

	"github.com/libp2p/go-libp2p/core/network"
)

func (wh *WingmanHandler) HandleStream(stream network.Stream) {
	defer stream.Close()

	for {
		decoder := json.NewDecoder(stream)
		var d WingmanMessage
		if err := decoder.Decode(&d); err != nil {
			wh.l.Printf("error handling stream: %v", err)
			return
		}

		wh.l.Println("received message:", d.Message)

		// if the message received has nothing to reconcile
		// continue

		// if the message received is missing one or more deletes
		// continue

		// if the message received has one or more deletes
		// delete the photos from the filesystem
		// reconcile CRDT

		// if the message received is missing one or more photos
		// send a message to this node with the photos it's missing

		// if the message received has one or more photos with all of the missing photos
		// add the photos to the filesystem
		// reconcile CRDT

		// if the message received has one or more photos with some of the missing photos
		// add the photos to the filesystem
		// reconcile CRDT with only the added photos
	}
}
