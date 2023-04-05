package wingman

import (
	"memory-lane/app/papaya"
	"memory-lane/app/raccoon"
)

type WingmanMessage struct {
	SenderMultiAddr string                    `json:"sender"`
	Album           string                    `json:"album"`
	Crdt            *raccoon.CRDT             `json:"crdt"`
	Photos          *map[string]*papaya.Photo `json:"photos,omitempty"`
}
