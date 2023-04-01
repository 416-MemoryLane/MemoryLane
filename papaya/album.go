package papaya

import "memory-lane/app/raccoon"

type Album struct {
	Crdt     *raccoon.CRDT
	Name     string
	Pictures Pictures
}

type Albums *map[raccoon.AlbumId]*Album
