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
func (g *Gallery) AddPicture() (*Picture, error) {
	return &Picture{}, nil
}
func (g *Gallery) DeletePicture() (*Picture, error) {
	return &Picture{}, nil
}
func (g *Gallery) GetAlbums() *Albums {
	return &Albums{}
}
func (g *Gallery) GetAllPicturesFromAlbum() *Pictures {
	return &Pictures{}
}
func (g *Gallery) GetPicture() *Picture {
	return &Picture{}
}
