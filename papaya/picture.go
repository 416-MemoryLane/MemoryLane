package papaya

import "fmt"

type Picture struct {
	Id  PictureId
	Img []byte
}

type PictureId string

type Pictures *map[PictureId]*Picture

var ErrPhotoExists = fmt.Errorf("picture exists")
