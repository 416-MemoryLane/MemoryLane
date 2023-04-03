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

		wh.l.Println("received message:", d)

		// if the gallery does not have the album, refetch from Galactus

		// Instantiate data structures required for comparing CRDTs
		// msgAlbum := d.Album
		// msgCrdt := d.Crdt
		// album := (*wh.Gallery.Albums)[msgAlbum]
		// albumCrdt := album.Crdt

		// Photos that are deleted and in the current node's album
		// delete the photos from the filesystem
		// reconcile CRDT

		// Photos that are deleted but that were never in the current node's album
		// reconcile CRDT by adding theses nodes to the deleted and added sets

		// Reconcile only the photos for which the actual images have been provided
		// add the photos to the filesystem
		// reconcile CRDT

		// If the current node has added photos that the sender node does not have
		// send a message to this node with the current CRDT
		// break

		// In the following cases, there is nothing to reconcile:
		// if the album states are equal
		// or if difference between album states is that the incoming album state is missing deletes
	}
}
