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
	AlbumName   string     `json:"album_name"`
	AddedList   *[]PhotoId `json:"added"`
	DeletedList *[]PhotoId `json:"deleted"`

	l *log.Logger
}

type AlbumId uuid.UUID

// Convert the AlbumId to string
func (aid AlbumId) String() string {
	return uuid.UUID(aid).String()
}

type PhotoId uuid.UUID

// Convert the PhotoId to string
func (aid PhotoId) String() string {
	return uuid.UUID(aid).String()
}

func NewCRDT(l *log.Logger) (*CRDT, error) {
	c := CRDT{
		&map[PhotoId]bool{},
		&map[PhotoId]bool{},

		NewAlbumId(),
		"",
		&[]PhotoId{},
		&[]PhotoId{},
		l,
	}

	return &c, nil
}

// Generate new id for a new album
func NewAlbumId() AlbumId {
	return AlbumId(uuid.New())
}

// Generate new id for a new photo
func NewPhotoId() PhotoId {
	return PhotoId(uuid.New())
}

// Return PhotoId from provided uuid string
func PhotoIdFromString(s string) PhotoId {
	return PhotoId(idFromString(s))
}

// Return AlbumId from provided uuid string
func AlbumIdFromString(s string) AlbumId {
	return AlbumId(idFromString(s))
}

func idFromString(s string) uuid.UUID {
	return uuid.Must(uuid.Parse(s))
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
