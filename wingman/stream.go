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

		// TODO: if the gallery does not have the album, refetch from Galactus

		// Instantiate data structures required for comparing CRDTs
		msgAlbumId := d.Album
		msgCrdt := d.Crdt
		album := (*wh.Gallery.Albums)[msgAlbumId]
		albumCrdt := album.Crdt

		// Find and handle photos to delete
		for p := range *msgCrdt.Deleted {
			msgVal, msgOk := (*msgCrdt.Deleted)[p]
			albumDeletedVal, albumDeletedOk := (*albumCrdt.Deleted)[p]

			if msgOk && msgVal && (!albumDeletedOk || !albumDeletedVal) {
				albumPhoto, err := wh.Gallery.GetPhoto(msgAlbumId, p)
				if err != nil {
					wh.l.Printf("error retrieving photo while reconciling node: %v\n", err)
					continue
				}

				// Reconcile file system and CRDT
				albumCrdt.AddPhoto(p)
				if albumPhoto != nil {
					wh.Gallery.DeletePhoto(msgAlbumId, p)
				}
				albumCrdt.DeletePhoto(p)
			}
		}

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
