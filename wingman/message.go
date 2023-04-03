package wingman

import "memory-lane/app/raccoon"

type WingmanMessage struct {
	SenderMultiAddr string              `json:"sender"`
	Album           string              `json:"album"`
	Crdt            *raccoon.CRDT       `json:"crdt"`
	Photos          *map[string]*[]byte `json:"photos,omitempty"`
}
