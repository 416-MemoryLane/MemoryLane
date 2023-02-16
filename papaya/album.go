package papaya

import "io/fs"

type Album struct {
	DirEntry fs.DirEntry
}

type Albums []Album
