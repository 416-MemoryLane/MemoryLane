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
	Albums raccoon.Albums
}

const GALLERY_DIR = "./memory-lane-gallery"

// Initialize a new gallery based on existing gallery in filesystem or create a new one if one doesn't exist
func NewGallery(l *log.Logger) (*Gallery, error) {
	gallery := &Gallery{l, &map[string]*raccoon.CRDT{}}

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

		albums := &map[string]*raccoon.CRDT{}
		gallery.Albums = albums
		for _, dir := range dirs {
			if dir.IsDir() {
				// Instantiate new album
				dirName := dir.Name()
				crdt := raccoon.NewCRDT(gallery.l)
				(*albums)[dirName] = crdt

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
						crdt = raccoon.NewCRDT(l)
						err = crdt.UnmarshalJSON(data)
						if err != nil {
							return nil, err
						}

						// Add crdt to albums
						(*albums)[dirName] = crdt

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
func (g *Gallery) CreateAlbum(albumName string) (*raccoon.CRDT, error) {
	// Initialise new CRDT with provided album name
	crdt := raccoon.NewCRDT(g.l)
	crdt.AlbumName = albumName

	// Create a new album directory
	dirName := crdt.Album
	albumDir := filepath.Join(GALLERY_DIR, dirName)
	err := os.Mkdir(albumDir, os.ModeDir|0777)
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

	// Add album to gallery
	(*g.Albums)[dirName] = crdt

	return crdt, nil
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
func (g *Gallery) GetAlbums() raccoon.Albums {
	return g.Albums
}

// Get an album. Return nil if the album does not exist
// TODO: This will likely be more complicated as it will be required to create the message to send to another node
func (g *Gallery) GetAlbum(aid string) *raccoon.CRDT {
	return (*g.Albums)[aid]
}

// Add a photo to an album if the album exists
func (g *Gallery) AddPhoto(aid string, photo []byte) (string, error) {
	crdt := (*g.Albums)[aid]
	if crdt == nil {
		return "", fmt.Errorf("album %s does not exist", aid)
	}

	// Decode the image bytes into an image
	p, _, err := image.Decode(bytes.NewReader(photo))
	if err != nil {
		return "", fmt.Errorf("failed to convert bytes to img: %w", err)
	}

	// Create a new file to save the image
	pid := uuid.New().String()
	// TODO: Must determine this based on image type
	photoFileName := fmt.Sprintf("%s.png", pid)
	photoFile := filepath.Join(GALLERY_DIR, aid, photoFileName)
	f, err := os.Create(photoFile)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	// Encode the photo to png and write it to the file
	if err := png.Encode(f, p); err != nil {
		return "", fmt.Errorf("failed to encode to png: %w", err)
	}

	// Add photo to CRDT and write to file
	crdt.AddPhoto(pid)
	crdtFile := filepath.Join(GALLERY_DIR, aid, "crdt.json")
	jsonData, err := crdt.MarshalJSON()
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON data: %w", err)
	}
	err = os.WriteFile(crdtFile, jsonData, 0777)
	if err != nil {
		return "", fmt.Errorf("failed to write file %s: %w", crdtFile, err)
	}

	return pid, nil
}

// Delete a photo from an album
func (g *Gallery) DeletePhoto(aid string, pid string) (string, error) {
	crdt := (*g.Albums)[aid]
	if crdt == nil {
		return "", fmt.Errorf("album %s does not exist", aid)
	}

	photo := (*crdt.Added)[pid]
	if !photo {
		return "", fmt.Errorf("photo %s does not exist", pid)
	}

	// Delete album from filesystem
	photoFile, err := filepath.Glob(filepath.Join(GALLERY_DIR, aid, pid+"*"))
	if err != nil {
		return "", fmt.Errorf("failed to find any photos matching %s: %w", photoFile, err)
	}
	err = os.Remove(photoFile[0])
	if err != nil {
		return "", fmt.Errorf("failed to delete photo %s: %w", photoFile, err)
	}

	// Remove photo from CRDT and write to file
	crdt.DeletePhoto(pid)
	crdtFile := filepath.Join(GALLERY_DIR, aid, "crdt.json")
	jsonData, err := crdt.MarshalJSON()
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON data: %w", err)
	}
	err = os.WriteFile(crdtFile, jsonData, 0777)
	if err != nil {
		return "", fmt.Errorf("failed to write file %s: %w", crdtFile, err)
	}

	return pid, nil
}

// Retrieve all the photo ids of an album
func (g *Gallery) GetPhotos(aid string) Photos {
	return (*g.Albums)[aid].Added
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
