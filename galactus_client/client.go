package galactus_client

import "log"

type GalactusClient struct {
	l *log.Logger
}

func NewGalactusClient(l *log.Logger) *GalactusClient {
	return &GalactusClient{l}
}
