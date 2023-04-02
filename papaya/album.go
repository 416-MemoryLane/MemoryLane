package papaya

import (
	"fmt"
	"memory-lane/app/raccoon"
)

type Album struct {
	Crdt   *raccoon.CRDT
	Photos Photos
}

type Albums *map[string]*Album

// Stringer for Album
func (a Album) String() string {
	return fmt.Sprintf("Album{CRDT: %v, Photos: %v}",
		*a.Crdt, *a.Photos)
}
