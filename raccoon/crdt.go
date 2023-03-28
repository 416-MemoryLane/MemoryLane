package raccoon

import "log"

type CRDT struct {
	// TODO: must add album when Papaya is added
	Added   *map[string]bool
	Deleted *map[string]bool

	l *log.Logger
}

func NewCRDT(l *log.Logger) *CRDT {
	return &CRDT{
		&map[string]bool{},
		&map[string]bool{},
		l,
	}
}

func (c *CRDT) AddPhoto(fn string) {
	(*c.Added)[fn] = true
}

func (c *CRDT) DeletePhoto(fn string) {
	delete(*c.Added, fn)
	(*c.Deleted)[fn] = true
}

func (c *CRDT) Reconcile(crdt *CRDT) (*CRDT, bool) {
	isChanged := false

	for k := range *crdt.Added {
		c.AddPhoto(k)
		isChanged = true
	}

	for k := range *crdt.Deleted {
		c.DeletePhoto(k)
		isChanged = true
	}

	return c, isChanged
}
