package papaya

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"memory-lane/app/raccoon"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type Gallery struct {
	l      *log.Logger
	Albums Albums
}

const GALLERY_DIR = "./memory-lane-gallery"

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

// Create a new album
func (g *Gallery) CreateAlbum(albumName string) (*Album, error) {
	// Initialise new CRDT with provided album name
	crdt, err := raccoon.NewCRDT(g.l)
	if err != nil {
		return nil, err
	}
	crdt.AlbumName = albumName

	// Create a new album directory
	dirName := crdt.Album
	albumDir := filepath.Join(GALLERY_DIR, dirName)
	err = os.Mkdir(albumDir, os.ModeDir|0777)
	if err != nil {
		return nil, fmt.Errorf("failed to create album: %w", err)
	}

	// Add CRDT to new album directory
	crdtFile := filepath.Join(albumDir, "crdt.json")
	jsonData, err := crdt.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON data: %w", err)
	}
	err = os.WriteFile(crdtFile, jsonData, 0777)
	if err != nil {
		return nil, fmt.Errorf("failed to write file %s: %w", crdtFile, err)
	}

	// Initialise a new album
	album := &Album{crdt, &map[string]bool{}}

	// Add album to gallery
	(*g.Albums)[dirName] = album

	return album, nil
}

// Delete an album if it exists
func (g *Gallery) DeleteAlbum(aid string) error {
	// Delete album from filesystem
	albumDir := filepath.Join(GALLERY_DIR, aid)
	err := os.RemoveAll(albumDir)
	if err != nil {
		return fmt.Errorf("failed to delete album %s: %w", albumDir, err)
	}

	// Remove from Gallery
	delete(*g.Albums, aid)

	return nil
}

// Retrieve all the albums (i.e. directories of photos) that the user is part of
func (g *Gallery) GetAlbums() Albums {
	return g.Albums
}

// Get an album. Return nil if the album does not exist
// TODO: This will likely be more complicated as it will be required to create the message to send to another node
func (g *Gallery) GetAlbum(aid string) *Album {
	return (*g.Albums)[aid]
}

// Add a photo to an album if the album exists
func (g *Gallery) AddPhoto(aid string, photo []byte) (string, error) {
	album := (*g.Albums)[aid]
	if album == nil {
		return "", fmt.Errorf("album %s does not exist", aid)
	}

	// Decode the image bytes into an image
	p, _, err := image.Decode(bytes.NewReader(photo))
	if err != nil {
		return "", fmt.Errorf("failed to convert bytes to img: %w", err)
	}

	// Create a new file to save the image
	pid := fmt.Sprintf("%s.png", uuid.New().String())
	photoFile := filepath.Join(GALLERY_DIR, aid, pid)
	f, err := os.Create(photoFile)
	if err != nil {
		return "", fmt.Errorf("failed to convert bytes to img: %w", err)
	}
	defer f.Close()

	// Encode the photo to png and write it to the file
	if err := png.Encode(f, p); err != nil {
		return "", fmt.Errorf("failed to encode to png: %w", err)
	}

	// Add photo to CRDT and write to file
	album.Crdt.AddPhoto(pid)
	crdtFile := filepath.Join(GALLERY_DIR, aid, "crdt.json")
	jsonData, err := album.Crdt.MarshalJSON()
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON data: %w", err)
	}
	err = os.WriteFile(crdtFile, jsonData, 0777)
	if err != nil {
		return "", fmt.Errorf("failed to write file %s: %w", crdtFile, err)
	}

	// Add photo to the Album
	(*album.Photos)[pid] = true

	return pid, nil
}

// Delete a photo from an album
func (g *Gallery) DeletePhoto(aid string, pid string) (string, error) {
	album := (*g.Albums)[aid]
	if album == nil {
		return "", fmt.Errorf("album %s does not exist", aid)
	}

	photo := (*album.Photos)[pid]
	if !photo {
		return "", fmt.Errorf("photo %s does not exist", pid)
	}

	// Delete album from filesystem
	photoFile := filepath.Join(GALLERY_DIR, aid, pid)
	err := os.Remove(photoFile)
	if err != nil {
		return "", fmt.Errorf("failed to delete photo %s: %w", photoFile, err)
	}

	// Remove photo from CRDT and write to file
	album.Crdt.DeletePhoto(pid)
	crdtFile := filepath.Join(GALLERY_DIR, aid, "crdt.json")
	jsonData, err := album.Crdt.MarshalJSON()
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON data: %w", err)
	}
	err = os.WriteFile(crdtFile, jsonData, 0777)
	if err != nil {
		return "", fmt.Errorf("failed to write file %s: %w", crdtFile, err)
	}

	// Remove photo from Album
	delete(*album.Photos, pid)

	return pid, nil
}

// Retrieve all the photo file names of an album
func (g *Gallery) GetPhotos(aid string) Photos {
	return (*g.Albums)[aid].Photos
}

// Retrieve the photo from an album
func (g *Gallery) GetPhoto(aid string, pid string) (*[]byte, error) {
	// Construct the file path to the photo based on the album ID and photo ID
	photoPath := filepath.Join(GALLERY_DIR, aid, pid)

	// Read the photo file into memory
	photoData, err := os.ReadFile(photoPath)
	if err != nil {
		// If there was an error reading the file, return it
		return nil, fmt.Errorf("error reading photo file: %v", err)
	}

	// Return a pointer to the photo data
	return &photoData, nil
}

// Stringer for Gallery
func (g Gallery) String() string {
	return fmt.Sprintf("Number of albums: %v, %v", len(*g.Albums), *g.Albums)
}
