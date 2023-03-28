package papaya

type Album struct {
	Id       AlbumId
	Name     string
	Pictures Pictures
}

type AlbumId string

type Albums *map[AlbumId]*Album
