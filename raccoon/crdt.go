package raccoon

import "log"

type CRDT struct {
	// TODO: must add album when Papaya is added
	added   *map[string]bool
	deleted *map[string]bool

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
	(*c.added)[fn] = true
}

func (c *CRDT) DeletePhoto(fn string) {
	delete(*c.added, fn)
	(*c.deleted)[fn] = true
}

func (c *CRDT) Reconcile(crdt *CRDT) *CRDT {
	for k := range *crdt.added {
		c.AddPhoto(k)
	}

	for k := range *crdt.deleted {
		c.DeletePhoto(k)
	}

	return c
}
