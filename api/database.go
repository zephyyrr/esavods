package main

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"io"
	"path"
)

const (
	eventsDBFile = "esavods_api_events.bolt"
	runsDBFile   = "esavods_api_runs.bolt"
)

var (
	db      *Database
	encoder MarshalUnmarshaller
)

type Database struct {
	Events, Users, Runs *bolt.DB
}

func OpenDatabase() (db *Database, err error) {
	db = &Database{}
	mo := multiopener{}
	db.Events = mo.Open(eventsDBFile)
	db.Runs = mo.Open(runsDBFile)
	err = mo.Error
	return
}

type multiopener struct {
	Error error
}

func (opener *multiopener) Open(file string) (db *bolt.DB) {
	db, err := bolt.Open(path.Join(DBFolder, file), 0600, nil)
	if err != nil && opener.Error == nil {
		opener.Error = err
	}
	return
}

type multicloser struct {
	Error error
}

func (closer *multicloser) Close(c io.Closer) {
	e := c.Close()
	if closer.Error == nil {
		closer.Error = e
	}
}

func (db *Database) Close() error {
	mc := multicloser{}
	mc.Close(db.Events)
	mc.Close(db.Runs)
	return mc.Error
}

type Boltable interface {
	ToBolt(*bolt.Bucket) error
	FromBolt(*bolt.Bucket) error
}

type MarshalUnmarshaller interface {
	Marshal(v interface{}) (data []byte, err error)
	Unmarshal(data []byte, v interface{}) error
}

// Json implementation of MarshalUnmarshaller
type Json struct{}

func (Json) Marshal(v interface{}) (data []byte, err error) {
	return json.Marshal(v)
}

func (Json) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
