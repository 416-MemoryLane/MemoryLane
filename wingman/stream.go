package wingman

import (
	"encoding/json"
	"memory-lane/app/papaya"
	"os"
	"path/filepath"

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
					_, err := wh.Gallery.DeletePhoto(msgAlbumId, p)
					if err != nil {
						wh.l.Printf("error deleting photo while reconciling node: %v\n", err)
						continue
					}
				}
				albumCrdt.DeletePhoto(p)
			}
		}

		// Find and handle photos to add
		msgPhotos := d.Photos
		if msgPhotos != nil {
			for p := range *msgPhotos {
				albumAddedVal, albumAddedOk := (*albumCrdt.Added)[p]

				// Reconcile file system and CRDT
				if !albumAddedOk || !albumAddedVal {
					_, err := wh.Gallery.AddPhoto(msgAlbumId, *(*msgPhotos)[p])
					if err != nil {
						wh.l.Printf("error adding photo while reconciling node: %v\n", err)
						continue
					}
				}
				albumCrdt.AddPhoto(p)
			}
		}

		// Create a message of photos to send to sender node
		var photosToSend map[string]*[]byte
		albumPhotos := wh.Gallery.GetPhotos(msgAlbumId)
		for p := range *albumPhotos {
			if val, ok := (*msgCrdt.Added)[p]; !val || !ok {
				d, err := wh.Gallery.GetPhoto(msgAlbumId, p)
				if err != nil {
					wh.l.Printf("error retrieving photo while creating reply message: %v\n", err)
					continue
				}
				photosToSend[p] = d
			}
		}

		// If there are photos to send, send WingmanMessage with the photos
		if photosToSend != nil {
			msg := WingmanMessage{
				SenderMultiAddr: wh.Multiaddr,
				Album:           msgAlbumId,
				Crdt:            albumCrdt,
				Photos:          &photosToSend,
			}

			encoder := json.NewEncoder(stream)
			if err := encoder.Encode(msg); err != nil {
				wh.l.Printf("error sending msg with added photos: %v\n", err)
			}
		}

		// Persist reconciled CRDT to filesystem
		crdtFile := filepath.Join(papaya.GALLERY_DIR, msgAlbumId, "crdt.json")
		jsonData, err := albumCrdt.MarshalJSON()
		if err != nil {
			wh.l.Printf("failed to marshal JSON data: %w", err)
		}
		err = os.WriteFile(crdtFile, jsonData, 0777)
		if err != nil {
			wh.l.Printf("failed to write file %s: %w", crdtFile, err)
		}

		// In the following cases, there is nothing to reconcile:
		// if the album states are equal
		// or if difference between album states is that the incoming album state is missing deletes
	}
}
