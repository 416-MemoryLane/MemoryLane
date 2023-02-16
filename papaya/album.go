package papaya

import "io/fs"

type Album struct {
	name     string
	fileMode fs.FileMode
	info     fs.FileInfo
}

type Albums []*Album

func (a Album) Name() string {
	return a.name
}

func (a Album) IsDir() bool {
	return true
}

func (a Album) Type() fs.FileMode {
	return a.fileMode
}

func (a Album) Info() (fs.FileInfo, error) {
	return a.info, nil
}
