package papaya

type Photo struct {
	ID       string  `json:"id"`
	MimeType string  `json:"mimetype"`
	Data     *[]byte `json:"data"`
}

type Photos *map[string]*Photo
