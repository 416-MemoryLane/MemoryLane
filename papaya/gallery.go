package papaya

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"memory-lane/app/raccoon"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type Gallery struct {
	GalleryDir string
	l          *log.Logger
}

// Initialize a new gallery based on existing gallery in filesystem or create a new one if one doesn't exist
func NewGallery(gd string, l *log.Logger) (*Gallery, error) {
	g := &Gallery{gd, l}

	// Try to open the gallery directory
	if _, err := os.Stat(g.GalleryDir); os.IsNotExist(err) {
		// If the gallery directory doesn't exist, create a new gallery directory and an empty album map
		err = os.Mkdir(g.GalleryDir, os.ModeDir|0777)
		if err != nil {
			return nil, fmt.Errorf("error instantiating a gallery: %v", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("error trying to access gallery: %v", err)
	}

	return g, nil
}

// Create a new album
func (g *Gallery) CreateAlbum(albumName string) (*raccoon.CRDT, error) {
	// Initialise new CRDT with provided album name and id
	albumId := uuid.New().String()
	crdt, err := g.AddAlbum(albumId, albumName)
	if err != nil {
		return nil, err
	}

	return crdt, nil
}

// Add a new album
func (g *Gallery) AddAlbum(albumId, albumName string) (*raccoon.CRDT, error) {
	// Initialise new CRDT with provided album name and id
	crdt := raccoon.NewCRDT(g.GalleryDir, albumId, albumName, g.l)

	// Create a new album directory
	albumDir := filepath.Join(g.GalleryDir, albumId)
	err := os.Mkdir(albumDir, os.ModeDir|0777)
	if err != nil {
		return nil, fmt.Errorf("failed to create album: %w", err)
	}

	// Add CRDT to new album directory
	err = crdt.PersistCRDT()
	if err != nil {
		return nil, fmt.Errorf("failed to persist CRDT of new album: %w", err)
	}

	return crdt, nil
}

// Delete an album if it exists
func (g *Gallery) DeleteAlbum(aid string) error {
	// Delete album from filesystem
	albumDir := filepath.Join(g.GalleryDir, aid)
	err := os.RemoveAll(albumDir)
	if err != nil {
		return fmt.Errorf("failed to delete album %s: %w", albumDir, err)
	}

	return nil
}

// Retrieve all the album IDs
func (g *Gallery) GetAlbumIDs() (*map[string]bool, error) {
	crdts, err := g.GetAlbumCRDTs()
	if err != nil {
		return nil, fmt.Errorf("failed to get album CRDTs while retrieving IDs")
	}

	ids := map[string]bool{}
	for _, crdt := range *crdts {
		ids[crdt.Album] = true
	}

	return &ids, nil
}

// Retrieve all the album CRDTs
func (g *Gallery) GetAlbumCRDTs() (raccoon.CRDTs, error) {
	// Read the directory contents
	albums, err := os.ReadDir(g.GalleryDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read gallery directory: %w", err)
	}

	// Loop through the albums and add to data structure of CRDTs
	crdts := map[string]*raccoon.CRDT{}
	for _, a := range albums {
		crdt, err := g.GetAlbumCRDT(a.Name())
		if err != nil {
			return nil, fmt.Errorf("error getting an album: %w", err)
		}

		crdts[crdt.Album] = crdt
	}

	return &crdts, nil
}

// Get an album's CRDT
func (g *Gallery) GetAlbumCRDT(aid string) (*raccoon.CRDT, error) {
	// Read data from filesystem
	crdtFile := filepath.Join(g.GalleryDir, aid, "crdt.json")
	crdtData, err := os.ReadFile(crdtFile)
	if err != nil {
		return nil, fmt.Errorf("error reading crdt file: %v", err)
	}

	// Deserialize data into CRDT struct
	crdt := &raccoon.CRDT{}
	crdt.GalleryDir = g.GalleryDir

	err = crdt.UnmarshalJSON(crdtData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling crdt data: %v", err)
	}

	return crdt, nil
}

// Add a photo to an album if the album exists
func (g *Gallery) AddPhotoWithFileName(aid, pid string, photo Photo) (string, error) {
	crdt, err := g.GetAlbumCRDT(aid)
	if err != nil {
		return "", fmt.Errorf("error getting albums info: %w", err)
	}

	// Register image formats
	image.RegisterFormat("jpg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)

	// Decode the image bytes
	p, _, err := image.Decode(bytes.NewReader(*photo.Data))
	if err != nil {
		return "", fmt.Errorf("failed to convert bytes to img: %w", err)
	}

	// Create and encode the photo based on mimeType
	switch mimeType := photo.MimeType; mimeType {
	case "image/jpg", "image/jpeg":
		suffix := "jpeg"
		photoFileName := fmt.Sprintf("%s.%s", pid, suffix)
		photoFile := filepath.Join(g.GalleryDir, aid, photoFileName)
		f, err := os.Create(photoFile)
		if err != nil {
			return "", fmt.Errorf("failed to create file: %w", err)
		}
		defer f.Close()

		if err := jpeg.Encode(f, p, nil); err != nil {
			return "", fmt.Errorf("failed to encode to %s: %w", suffix, err)
		}
	case "image/png":
		suffix := "png"
		photoFileName := fmt.Sprintf("%s.%s", pid, suffix)
		photoFile := filepath.Join(g.GalleryDir, aid, photoFileName)
		f, err := os.Create(photoFile)
		if err != nil {
			return "", fmt.Errorf("failed to create file: %w", err)
		}
		defer f.Close()

		pngEncoder := png.Encoder{
			CompressionLevel: png.BestCompression,
		}
		if err := pngEncoder.Encode(f, p); err != nil {
			return "", fmt.Errorf("failed to encode to %s: %w", suffix, err)
		}
	default:
		return "", fmt.Errorf("failed to convert bytes to img: unsupported type")
	}

	// Add photo to CRDT and persist to file
	err = crdt.AddPhoto(pid)
	if err != nil {
		return "", fmt.Errorf("failed to persist CRDT: %w", err)
	}

	return pid, nil
}

// Delete a photo from an album
func (g *Gallery) DeletePhoto(aid string, pid string) (string, error) {
	crdt, err := g.GetAlbumCRDT(aid)
	if err != nil {
		return "", fmt.Errorf("error retrieving album crdt %w", err)
	}

	// Retrieve photo from filesystem
	_, err = g.GetPhoto(aid, pid)
	if err != nil {
		if os.IsNotExist(err) {
			// If the photo to be deleted does not exist, do nothing
			return "", nil
		} else {
			return "", fmt.Errorf("error retrieving photo: %w", err)
		}
	}

	// Delete photo from filesystem
	photoFile, err := filepath.Glob(filepath.Join(g.GalleryDir, aid, pid+"*"))
	if err != nil {
		return "", fmt.Errorf("failed to find any photos matching %s: %w", photoFile, err)
	}
	err = os.Remove(photoFile[0])
	if err != nil {
		return "", fmt.Errorf("failed to delete photo %s: %w", photoFile, err)
	}

	// Remove photo from CRDT and write to file
	err = crdt.DeletePhoto(pid)
	if err != nil {
		return "", fmt.Errorf("failed to delete photo in CRDT: %w", err)
	}

	return pid, nil
}

// Retrieve all the photos of an album
func (g *Gallery) GetPhotos(aid string) (Photos, error) {
	// Get the current album's CRDT
	crdt, err := g.GetAlbumCRDT(aid)
	if err != nil {
		return nil, fmt.Errorf("error retrieving album crdt: %w", err)
	}

	// Add all photos to Photos struct
	photos := map[string]*Photo{}
	for pid := range *crdt.Added {
		photo, err := g.GetPhoto(aid, pid)
		if err != nil {
			return nil, fmt.Errorf("error retrieving photos: %w", err)
		}

		photos[pid] = photo
	}

	return &photos, nil
}

// Retrieve the photo from an album
func (g *Gallery) GetPhoto(aid string, pid string) (*Photo, error) {
	// Construct the file path to the photo based on the album ID and photo ID
	photoPath, err := filepath.Glob(filepath.Join(g.GalleryDir, aid, pid+"*"))
	if err != nil || len(photoPath) == 0 {
		return nil, fmt.Errorf("failed to find any photos matching %s: %w", photoPath, err)
	}

	// Check if photo exists
	_, err = os.Stat(photoPath[0])
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist
			return nil, err
		}
		// Other error occurred, return it
		return nil, err
	}

	// Read the photo file into memory
	photoData, err := os.ReadFile(photoPath[0])
	if err != nil {
		// If there was an error reading the file, return it
		return nil, fmt.Errorf("error reading photo file: %v", err)
	}

	// Determine the MIME type based on the photo file contents
	mimeType := http.DetectContentType(photoData)

	// Return a pointer to the photo data
	return &Photo{pid, mimeType, &photoData}, nil
}
