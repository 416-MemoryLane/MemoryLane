package wingman

import "memory-lane/app/raccoon"

type WingmanMessage struct {
	SenderMultiAddr string        `json:"sender"`
	Album           string        `json:"album"`
	Crdt            *raccoon.CRDT `json:"crdt"`
	Photos          *[][]byte     `json:"photos;omitifempty"`
}
