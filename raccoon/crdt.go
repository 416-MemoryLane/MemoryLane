package raccoon

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

type CRDT struct {
	Added   *map[PhotoId]bool `json:"-"`
	Deleted *map[PhotoId]bool `json:"-"`

	Album       AlbumId    `json:"album"`
	AddedList   *[]PhotoId `json:"added"`
	DeletedList *[]PhotoId `json:"deleted"`

	l *log.Logger
}

type AlbumId uuid.UUID
type PhotoId uuid.UUID

func NewCRDT(l *log.Logger) *CRDT {
	return &CRDT{
		&map[PhotoId]bool{},
		&map[PhotoId]bool{},

		NewAlbumId(),
		&[]PhotoId{},
		&[]PhotoId{},
		l,
	}
}

// Generate new id for a new album
func NewAlbumId() AlbumId {
	return AlbumId(uuid.New())
}

// Generate new id for a new photo
func NewPhotoId() PhotoId {
	return PhotoId(uuid.New())
}

func (c *CRDT) AddPhoto(p PhotoId) {
	(*c.Added)[p] = true
}

func (c *CRDT) DeletePhoto(p PhotoId) {
	delete(*c.Added, p)
	(*c.Deleted)[p] = true
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
		AddedList   *[]PhotoId `json:"added"`
		DeletedList *[]PhotoId `json:"deleted"`
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
		AddedList   *[]PhotoId `json:"added"`
		DeletedList *[]PhotoId `json:"deleted"`
	}{
		CRDTAlias:   (*CRDTAlias)(c),
		AddedList:   listFromMap(c.Added),
		DeletedList: listFromMap(c.Deleted),
	})
}

func listFromMap(m *map[PhotoId]bool) *[]PhotoId {
	l := make([]PhotoId, 0, len(*m))
	for k := range *m {
		l = append(l, k)
	}
	return &l
}

func mapFromList(l *[]PhotoId) *map[PhotoId]bool {
	m := make(map[PhotoId]bool, len(*l))
	for _, k := range *l {
		m[k] = true
	}
	return &m
}
