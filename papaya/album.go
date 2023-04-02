package papaya

import "memory-lane/app/raccoon"

type Album struct {
	Crdt   *raccoon.CRDT
	Photos Photos
}

type Albums *map[raccoon.AlbumId]*Album
