package raccoon

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// TODO: Add as a field in the CRDT struct (`json:-`)
const GALLERY_DIR = "./memory-lane-gallery"

type CRDT struct {
	Added   *map[string]bool `json:"-"`
	Deleted *map[string]bool `json:"-"`

	Album       string    `json:"album"`
	AlbumName   string    `json:"album_name"`
	AddedList   *[]string `json:"added"`
	DeletedList *[]string `json:"deleted"`

	l *log.Logger
}

type CRDTs *map[string]*CRDT

func NewCRDT(albumId, albumName string, l *log.Logger) *CRDT {
	c := CRDT{
		&map[string]bool{},
		&map[string]bool{},

		albumId,
		albumName,
		&[]string{},
		&[]string{},
		l,
	}

	return &c
}

// Write CRDT data into the filesystem
func (c *CRDT) PersistCRDT() error {
	crdtFile := filepath.Join(GALLERY_DIR, c.Album, "crdt.json")
	jsonData, err := c.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal CRDT JSON data: %w", err)
	}
	err = os.WriteFile(crdtFile, jsonData, 0777)
	if err != nil {
		return fmt.Errorf("failed to write CRDT file: %w", err)
	}

	return nil
}

// Add photo to CRDT and persist to filesystem
func (c *CRDT) AddPhoto(pid string) error {
	(*c.Added)[pid] = true

	err := c.PersistCRDT()
	if err != nil {
		delete(*c.Added, pid)
		return err
	}

	return nil
}

// Add deleted photo to CRDT and persist to filesystem
func (c *CRDT) DeletePhoto(pid string) error {
	delete(*c.Added, pid)
	(*c.Deleted)[pid] = true

	err := c.PersistCRDT()
	if err != nil {
		delete(*c.Deleted, pid)
		(*c.Added)[pid] = true
		return err
	}

	return nil
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
