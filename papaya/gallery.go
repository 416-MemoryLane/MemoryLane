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

		err = os.Mkdir(relPath, fs.ModeDir)
		if err != nil {
			return nil, err
		}
	}

	var a Albums
	for _, e := range d {
		i, err := e.Info()
		if err != nil {
			return nil, err
		}

		a = append(a, &Album{e.Name(), e.Type(), i})
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

	err := os.Mkdir(albumPath, os.ModeDir)
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

	g.Albums = append(g.Albums, &Album{e.Name(), e.Type(), info})

	return albumPath, nil
}

// Stringer for Gallery
func (g Gallery) String() string {
	return fmt.Sprintf("\nGallery filepath: %v\nNumber of albums: %v", g.FullPath, len(g.Albums))
}
