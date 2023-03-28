package papaya

import (
	"fmt"
	"log"
)

type Gallery struct {
	l      *log.Logger
	Albums Albums
}

func NewGallery(l *log.Logger) (*Gallery, error) {
	return &Gallery{l, nil}, nil
}

// Create a new album if it doesn't exist
func (g *Gallery) CreateAlbum(aid string) (string, error) {
	// Must initialise a new CRDT and add it to the filesystem
	// Must create a new directory for its photos
	// Must create a new entry for Galactus

	return "", nil
}

// Delete an album
func (g *Gallery) DeleteAlbum(aid string) (string, error) {
	// Must delete the chosen directory and all its contents (photos + CRDT)
	// Must delete the entry in Galactus

	return "", nil
}

// Retrieve all the albums (i.e. directories of photos) that the user is part of
func (g *Gallery) GetAlbums() (Albums, error) {
	return nil, nil
}

// Get an album
func (g *Gallery) GetAlbum(album string) (Album, error) {
	return Album{}, nil
}

// Add a picture to an album
func (g *Gallery) AddPicture(aid string, picture Picture) (Picture, error) {
	// Must also update the CRDT in the filesystem

	return Picture{}, nil
}

// Delete a picture from an album
func (g *Gallery) DeletePicture(aid string, pid string) (Picture, error) {
	// Must also update the CRDT in the filesystem

	return Picture{}, nil
}

// Retrieve all the photos of an album
func (g *Gallery) GetPictures() (Pictures, error) {
	return nil, nil
}

// Retrieve the photo from an album
func (g *Gallery) GetPicture(aid string, pid string) (Picture, error) {
	return Picture{}, nil
}

// Stringer for Gallery
func (g Gallery) String() string {
	return fmt.Sprintf("\nNumber of albums: %v", len(*g.Albums))
}
