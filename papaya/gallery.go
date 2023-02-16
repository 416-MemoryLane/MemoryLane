package papaya

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
)

type Gallery struct {
	l        *log.Logger
	FullPath string
	RelPath  string
	Albums   Albums
}

func NewGallery(l *log.Logger, relPath string) (*Gallery, error) {
	d, err := os.ReadDir(relPath)

	if os.IsNotExist(err) {
		l.Printf("Gallery at %s does not exist", relPath)
		l.Printf("Creating a new gallery at %s", relPath)

		err = os.Mkdir(relPath, 0777)
		if err != nil {
			return nil, err
		}
	}

	var a Albums = &map[string]*Album{}
	for _, e := range d {
		var p Photos = &map[string]bool{}
		(*a)[e.Name()] = &Album{e.Name(), p}
	}

	wd, err := os.Getwd()
	if err != nil {
		l.Println(err)
	}
	fullPath := fmt.Sprintf("%v/%v", wd, relPath)
	relPath = fmt.Sprintf("./%v", relPath)

	return &Gallery{l, fullPath, relPath, a}, nil
}

// Create an album if it doesn't exist
func (g *Gallery) CreateAlbum(album string) (string, error) {
	albumPath := fmt.Sprintf("%v/%v", g.RelPath, album)

	err := os.Mkdir(albumPath, 0777)
	if err != nil {
		return "", err
	}
	g.l.Printf("New album created: %v", album)

	i, err := os.Stat(albumPath)
	if err != nil {
		return "", err
	}
	e := fs.FileInfoToDirEntry(i)

	var p Photos = &map[string]bool{}
	(*g.Albums)[e.Name()] = &Album{e.Name(), p}

	return albumPath, nil
}

// Delete an album
func (g *Gallery) DeleteAlbum(album string) (string, error) {
	albumPath := fmt.Sprintf("%v/%v", g.RelPath, album)

	_, err := os.Stat(albumPath)
	if err != nil {
		return "", err
	}

	err = os.Remove(albumPath)
	if err != nil {
		return "", err
	}

	delete(*g.Albums, album)

	g.l.Printf("Album deleted: %v", album)

	return albumPath, err
}

// Add a photo to an album
func (g *Gallery) AddPhoto(album string, photo Photo) (Photo, error) {
	albumPath := fmt.Sprintf("%v/%v", g.RelPath, album)

	_, err := os.Stat(albumPath)
	if err != nil {
		return "", err
	}

	url := "https://i.imgur.com/phzB4jR.png"
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	photoPath := fmt.Sprintf("%v/%v.png", albumPath, photo)
	_, err = os.Stat(photoPath)
	if err == nil {
		return "", ErrPhotoExists
	} else if !os.IsNotExist(err) {
		return "", err
	}

	f, err := os.Create(photoPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return "", err
	}

	g.l.Printf("New photo added to %v: %v", album, photo)

	return photo, nil
}

// Stringer for Gallery
func (g Gallery) String() string {
	return fmt.Sprintf("\nGallery filepath: %v\nNumber of albums: %v", g.FullPath, len(*g.Albums))
}
