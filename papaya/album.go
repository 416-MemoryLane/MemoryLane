package papaya

type Album struct{}

type Albums []*Album

func (g *Gallery) AddPicture() (*Picture, error) {
	return &Picture{}, nil
}

func (g *Gallery) DeletePicture() (*Picture, error) {
	return &Picture{}, nil
}

func (g *Gallery) GetPicture() *Picture {
	return &Picture{}
}

func (g *Gallery) GetAllPictures() *Pictures {
	return &Pictures{}
}
