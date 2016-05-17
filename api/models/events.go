package models

import (
	"github.com/boltdb/bolt"
	"net/http"
	"time"
)

// Represents a group of events.
type Events []Event

// A single event and/or marathon
type Event struct {
	Name  string `json:"name"`
	Dates Dates  `json:"dates"`
}

func (e Event) ToBolt(bucket *bolt.Bucket) error {
	// The Puts can not produce errors since they are bounded operations with static keys and limited encoding lengths in a writable transaction.
	bucket.Put([]byte("Name"), []byte(e.Name))
	dates, _ := bucket.CreateBucketIfNotExists([]byte("Dates"))
	return e.Dates.ToBolt(dates)
}

func (e *Event) FromBolt(bucket *bolt.Bucket) error {
	e.Name = string(bucket.Get([]byte("Name")))
	dates := bucket.Bucket([]byte("Dates"))
	return e.Dates.FromBolt(dates)
}

type Dates struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func (d Dates) ToBolt(bucket *bolt.Bucket) error {
	bucket.Put([]byte("Start"), []byte(d.Start.Format(time.RFC1123Z)))
	bucket.Put([]byte("End"), []byte(d.End.Format(time.RFC1123Z)))
	return nil
}

func (d *Dates) FromBolt(bucket *bolt.Bucket) (err error) {
	d.Start, err = time.Parse(time.RFC1123Z, string(bucket.Get([]byte("Start"))))
	if err != nil {
		return Error{HttpStatus: http.StatusInternalServerError, Message: "Unable to parse start date.", Internal: err}
	}
	d.End, err = time.Parse(time.RFC1123Z, string(bucket.Get([]byte("End"))))
	if err != nil {
		return Error{HttpStatus: http.StatusInternalServerError, Message: "Unable to parse end date.", Internal: err}
	}
	return nil
}
