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

	l.Printf("%v", g)

	// Most likely will want to reconcile gallery with other nodes here
}
