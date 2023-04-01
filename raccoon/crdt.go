package raccoon

import (
	"encoding/json"
	"log"
)

type CRDT struct {
	Added   *map[string]bool `json:"-"`
	Deleted *map[string]bool `json:"-"`

	AddedList   *[]string `json:"added"`
	DeletedList *[]string `json:"deleted"`

	l *log.Logger
}

func NewCRDT(l *log.Logger) *CRDT {
	return &CRDT{
		&map[string]bool{},
		&map[string]bool{},
		&[]string{},
		&[]string{},
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

func (c *CRDT) UnmarshalJSON(d []byte) error {
	type CRDTAlias CRDT
	aux := &struct {
		*CRDTAlias
		AddedList   *[]string `json:"added"`
		DeletedList *[]string `json:"deleted"`
	}{
		CRDTAlias: (*CRDTAlias)(c),
	}

	if err := json.Unmarshal(d, &aux); err != nil {
		return err
	}

	c.Added = mapFromList(aux.AddedList)
	c.Deleted = mapFromList(aux.DeletedList)

	return nil
}

func (c *CRDT) MarshalJSON() ([]byte, error) {
	type CRDTAlias CRDT
	return json.Marshal(&struct {
		*CRDTAlias
		AddedList   *[]string `json:"added"`
		DeletedList *[]string `json:"deleted"`
	}{
		CRDTAlias:   (*CRDTAlias)(c),
		AddedList:   listFromMap(c.Added),
		DeletedList: listFromMap(c.Deleted),
	})
}

func listFromMap(m *map[string]bool) *[]string {
	l := make([]string, 0, len(*m))
	for k := range *m {
		l = append(l, k)
	}
	return &l
}

func mapFromList(l *[]string) *map[string]bool {
	m := make(map[string]bool, len(*l))
	for _, k := range *l {
		m[k] = true
	}
	return &m
}
