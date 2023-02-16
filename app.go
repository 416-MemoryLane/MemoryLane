package main

import (
	"log"
	"memory-lane/app/papaya"
	"os"
)

func main() {
	GALLERY_REL_PATH := "gallery"

	l := log.New(os.Stdout, "memory-lane ", log.LstdFlags)
	g, err := papaya.NewGallery(l, GALLERY_REL_PATH)
	if err != nil {
		l.Fatal("Error initializing gallery: ", err)
	}

	// a, err := g.CreateAlbum("cool-album4")
	// if err != nil {
	// 	l.Println("Error creating album:", err)
	// } else {
	// 	l.Printf(a)
	// }

	// a, err := g.DeleteAlbum("gallery")
	// if err != nil {
	// 	l.Println("Error deleting album:", err)
	// } else {
	// 	l.Printf(a)
	// }

	l.Printf("%v", g)

	// Most likely will want to reconcile gallery with other nodes here
}
