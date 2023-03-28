package papaya

import "memory-lane/app/raccoon"

type Album struct {
	Id       AlbumId
	Crdt     *raccoon.CRDT
	Name     string
	Pictures Pictures
}

type AlbumId string

type Albums *map[AlbumId]*Album
