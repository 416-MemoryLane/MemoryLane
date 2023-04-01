package raccoon

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

type CRDT struct {
	Added   *map[uuid.UUID]bool `json:"-"`
	Deleted *map[uuid.UUID]bool `json:"-"`

	Id          uuid.UUID    `json:"id"`
	AddedList   *[]uuid.UUID `json:"added"`
	DeletedList *[]uuid.UUID `json:"deleted"`

	l *log.Logger
}

func NewCRDT(l *log.Logger) *CRDT {
	// Generate a new UUID
	id := uuid.New()

	return &CRDT{
		&map[uuid.UUID]bool{},
		&map[uuid.UUID]bool{},

		id,
		&[]uuid.UUID{},
		&[]uuid.UUID{},
		l,
	}
}

func (c *CRDT) AddPhoto(fn uuid.UUID) {
	(*c.Added)[fn] = true
}

func (c *CRDT) DeletePhoto(fn uuid.UUID) {
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
		AddedList   *[]uuid.UUID `json:"added"`
		DeletedList *[]uuid.UUID `json:"deleted"`
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
		AddedList   *[]uuid.UUID `json:"added"`
		DeletedList *[]uuid.UUID `json:"deleted"`
	}{
		CRDTAlias:   (*CRDTAlias)(c),
		AddedList:   listFromMap(c.Added),
		DeletedList: listFromMap(c.Deleted),
	})
}

func listFromMap(m *map[uuid.UUID]bool) *[]uuid.UUID {
	l := make([]uuid.UUID, 0, len(*m))
	for k := range *m {
		l = append(l, k)
	}
	return &l
}

func mapFromList(l *[]uuid.UUID) *map[uuid.UUID]bool {
	m := make(map[uuid.UUID]bool, len(*l))
	for _, k := range *l {
		m[k] = true
	}
	return &m
}
