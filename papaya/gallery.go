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

// Stringer for Gallery
func (g Gallery) String() string {
	return fmt.Sprintf("\nGallery filepath: %v\nNumber of albums: %v", g.FullPath, len(g.Albums))
}
