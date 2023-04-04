package raccoon

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type CRDT struct {
	Added   *map[string]bool `json:"-"`
	Deleted *map[string]bool `json:"-"`

	Album       string    `json:"album"`
	AlbumName   string    `json:"album_name"`
	AddedList   *[]string `json:"added"`
	DeletedList *[]string `json:"deleted"`

	l *log.Logger
}

type Albums *map[string]*CRDT

func NewCRDT(l *log.Logger) *CRDT {
	c := CRDT{
		&map[string]bool{},
		&map[string]bool{},

		uuid.New().String(),
		"",
		&[]string{},
		&[]string{},
		l,
	}

	return &c
}

func (c *CRDT) AddPhoto(pid string) {
	(*c.Added)[pid] = true
}

func (c *CRDT) DeletePhoto(pid string) {
	delete(*c.Added, pid)
	(*c.Deleted)[pid] = true
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

// Stringer for CRDT
func (c CRDT) String() string {
	return fmt.Sprintf("CRDT{Album: %s, AlbumName: %s, Added: %v, Deleted: %v}",
		c.Album, c.AlbumName, *c.Added, *c.Deleted)
}
