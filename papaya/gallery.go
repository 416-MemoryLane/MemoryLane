package papaya

import (
	"fmt"
	"io/fs"
	"log"
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
		i, err := e.Info()
		if err != nil {
			return nil, err
		}

		(*a)[e.Name()] = &Album{e.Name(), e.Type(), i}
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
	info, err := e.Info()
	if err != nil {
		return "", err
	}

	(*g.Albums)[e.Name()] = &Album{e.Name(), e.Type(), info}

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

// Stringer for Gallery
func (g Gallery) String() string {
	return fmt.Sprintf("\nGallery filepath: %v\nNumber of albums: %v", g.FullPath, len(*g.Albums))
}
