package papaya

import (
	"fmt"
	"log"
	"memory-lane/app/raccoon"
	"os"
	"path/filepath"
)

type Gallery struct {
	l      *log.Logger
	Albums Albums
}

// Initialize a new gallery based on existing gallery in filesystem or create a new one if one doesn't exist
func NewGallery(l *log.Logger) (*Gallery, error) {
	gallery := &Gallery{l, &map[string]*Album{}}

	// Define the path to the gallery directory
	galleryDir := "./memory-lane-gallery"

	// Try to open the gallery directory
	_, err := os.Stat(galleryDir)
	if err == nil {
		// If the directory exists, use the contents of the gallery to instantiate a new gallery
		// for each album, instantiate an album and CRDT
		d, err := os.Open(galleryDir)
		if err != nil {
			return nil, err
		}
		defer d.Close()

		// Read all directories
		dirs, err := d.Readdir(-1)
		if err != nil {
			return nil, err
		}

		albums := &map[string]*Album{}
		gallery.Albums = albums
		for _, dir := range dirs {
			if dir.IsDir() {
				// Instantiate new album
				dirName := dir.Name()
				crdt := &raccoon.CRDT{}
				photos := &map[string]bool{}
				album := &Album{crdt, photos}
				(*albums)[dirName] = album

				// Read all photos and CRDT to use to instantiate
				albumDir := filepath.Join(galleryDir, dirName)
				files, err := os.ReadDir(albumDir)
				if err != nil {
					return nil, err
				}

				for _, file := range files {
					fileName := file.Name()
					filePath := filepath.Join(albumDir, fileName)

					// Read the contents of the file
					data, err := os.ReadFile(filePath)
					if err != nil {
						return nil, err
					}

					if fileName == "crdt.json" {
						// If crdt.json, deserialize to CRDT struct
						crdt, err = raccoon.NewCRDT(l)
						if err != nil {
							return nil, err
						}
						err = crdt.UnmarshalJSON(data)
						if err != nil {
							return nil, err
						}

						// Must add to album.Crdt here even though its pointer has been added to the album
						album.Crdt = crdt

					} else {
						// Else add to album's map of photos
						(*photos)[fileName] = true
					}
				}
			}
		}

		return gallery, nil
	}

	// If the directory doesn't exist, create a new gallery directory and an empty album map
	err = os.Mkdir(galleryDir, os.ModeDir|0777)
	if err != nil {
		return nil, err
	}

	return gallery, nil
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

// Add a photo to an album
// TODO: fix return value
func (g *Gallery) AddPhoto(aid string, photo []byte) (interface{}, error) {
	// Must also update the CRDT in the filesystem

	return nil, nil
}

// Delete a photo from an album
// TODO: fix return value
func (g *Gallery) DeletePhoto(aid string, pid string) (interface{}, error) {
	// Must also update the CRDT in the filesystem

	return nil, nil
}

// Retrieve all the photos of an album
// TODO: fix return value
func (g *Gallery) GetPhotos() (interface{}, error) {
	return nil, nil
}

// Retrieve the photo from an album
// TODO: fix return value
func (g *Gallery) GetPhoto(aid string, pid string) (interface{}, error) {
	return nil, nil
}

// Stringer for Gallery
func (g Gallery) String() string {
	return fmt.Sprintf("Number of albums: %v, %v", len(*g.Albums), *g.Albums)
}
