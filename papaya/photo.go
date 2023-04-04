package papaya

type Photo struct {
	MimeType string  `json:"mimetype"`
	Data     *[]byte `json:"data"`
}

type Photos *map[string]bool
