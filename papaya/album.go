package papaya

type Album struct {
	Name   string
	Photos Photos
}

type Albums *map[string]*Album
