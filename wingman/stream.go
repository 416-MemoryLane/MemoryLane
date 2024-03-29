package wingman

import (
	"context"
	"encoding/json"
	"memory-lane/app/papaya"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
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

		// Retrieve CRDT from filesystem
		albumCrdt, err := wh.Gallery.GetAlbumCRDT(msgAlbumId)
		if err != nil {
			wh.l.Printf("error retrieving album's crdt: %v\n", err)
			continue
		}

		// Find and handle photos to delete
		for p := range *msgCrdt.Deleted {
			msgVal, msgOk := (*msgCrdt.Deleted)[p]
			albumDeletedVal, albumDeletedOk := (*albumCrdt.Deleted)[p]

			// If there is a deleted photo in the sender node but not the current node
			if msgOk && msgVal && (!albumDeletedOk || !albumDeletedVal) {
				// Reconcile file system and CRDT
				// TODO: wh.Gallery.DeletePhoto and albumCrdt.DeletePhoto should be atomic
				_, err := wh.Gallery.DeletePhoto(msgAlbumId, p)
				if err != nil {
					wh.l.Printf("error deleting photo while reconciling node: %v\n", err)
				}
				err = albumCrdt.DeletePhoto(p)
				if err != nil {
					wh.l.Printf("error reconciling CRDT in a delete operation: %v\n", err)
					continue
				}
			}
		}

		// Find and handle photos to add
		msgPhotos := d.Photos
		if msgPhotos != nil {
			for p := range *msgPhotos {
				albumAddedVal, albumAddedOk := (*albumCrdt.Added)[p]

				// Reconcile file system and CRDT
				// TODO: albumCrdt.AddPhoto and wh.Gallery.AddPhotoWithFileName should be atomic
				if !albumAddedOk || !albumAddedVal {
					_, err := wh.Gallery.AddPhotoWithFileName(msgAlbumId, p, *(*msgPhotos)[p])
					if err != nil {
						wh.l.Printf("error adding photo while reconciling node: %v\n", err)
						continue
					}
					err = albumCrdt.AddPhoto(p)
					if err != nil {
						wh.l.Printf("error reconciling CRDT in a add operation: %v\n", err)
						continue
					}
				}
			}
		}

		// Create a message of photos to send to sender node
		var photosToSend map[string]*papaya.Photo
		albumPhotos, err := wh.Gallery.GetPhotos(msgAlbumId)
		if err != nil {
			wh.l.Printf("error getting photos: %v\n", err)
			continue
		}
		for p := range *albumPhotos {
			if val, ok := (*msgCrdt.Added)[p]; !val || !ok {
				if photosToSend == nil {
					photosToSend = make(map[string]*papaya.Photo)
				}
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

			// Connect to the node at the given address
			peerAddrInfo := wh.MultiaddrStrToPeerAddrInfo(d.SenderMultiAddr)
			if err := (*wh.Node).Connect(context.Background(), *peerAddrInfo); err != nil {
				panic(err)
			}
			wh.l.Println("connected to in handler:", peerAddrInfo.String())

			// Open a new stream to a connected node
			s, err := (*wh.Node).NewStream(context.Background(), peerAddrInfo.ID, protocol.ID(wh.ProtocolId))
			if err != nil {
				wh.l.Printf("failed opening a new stream: %v", err)
			}

			encoder := json.NewEncoder(s)
			if err := encoder.Encode(msg); err != nil {
				wh.l.Printf("error sending msg with added photos: %v\n", err)
			}

			wh.l.Printf("sent msg to: %v\n for album: %v from handler\n", wh.Multiaddr, msgAlbumId)
		}

		// In the following cases, there is nothing to reconcile:
		// if the album states are equal
		// or if difference between album states is that the incoming album state is missing deletes
	}
}
