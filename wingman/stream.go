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
		msgAlbum := d.Album
		msgCrdt := d.Crdt
		album := (*wh.Gallery.Albums)[msgAlbum]
		albumCrdt := album.Crdt

		for {
			// Album states are equal so there is nothing to reconcile
			if albumCrdt.Equals(msgCrdt) {
				break
			}

			// if the message received is missing one or more deletes
			// continue

			// if the message received has one or more deletes
			// delete the photos from the filesystem
			// reconcile CRDT

			// if the message received is missing one or more photos
			// AND these photos have not already been deleted
			// send a message to this node with the photos it's missing

			// if the message received has one or more photos with all of the missing photos
			// add the photos to the filesystem
			// reconcile CRDT

			// if the message received has one or more photos with some of the missing photos
			// add the photos to the filesystem
			// reconcile CRDT with only the added photos
		}
	}
}
