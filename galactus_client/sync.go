package galactus_client

type SyncAlbum struct {
	AlbumID         string   `json:"albumId"`
	AlbumName       string   `json:"albumName"`
	AuthorizedUsers []string `json:"authorizedUsers"`
	CreatedBy       string   `json:"createdBy"`
}

type SyncResponse []*SyncAlbum
