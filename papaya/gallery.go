package papaya

import (
	"log"
)

type Gallery struct {
	l        *log.Logger
	filepath string
}

func NewGallery(l *log.Logger, filepath string) *Gallery {
	return &Gallery{l, filepath}
}

func (g *Gallery) CreateAlbum() (*Album, error) {
	return &Album{}, nil
}

func (g *Gallery) DeleteAlbum() (*Album, error) {
	return &Album{}, nil
}

func (g *Gallery) GetAllAlbums() *Albums {
	return &Albums{}
}
