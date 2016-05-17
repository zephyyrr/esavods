package main

import (
	"github.com/boltdb/bolt"
	"github.com/labstack/echo"
	. "github.com/zephyyrr/esavods/api/models"
	"net/http"
)

func GetEvents(c echo.Context) error {
	return db.Events.View(func(tx *bolt.Tx) error {
		events := make([]Event, 0, 16) //Preallocate space for 16 events.

		err := tx.ForEach(func(name []byte, b *bolt.Bucket) (err error) {
			var event Event
			err = event.FromBolt(b)
			log.WithField("error", err).WithField("event", event).Debug("Fetched event")
			events = append(events, event)
			return
		})
		if err != nil {

		}
		return c.JSON(http.StatusOK, events)
	})
}

func GetEvent(c echo.Context) error {
	return db.Events.View(func(tx *bolt.Tx) error {
		var event Event
		b := tx.Bucket([]byte(c.P(0)))
		if b == nil {
			return Error{HttpStatus: http.StatusNotFound, Message: "The requested event was not found.", Data: c.P(0)}
		}
		err := event.FromBolt(b)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, event)
	})
}

func PostEvents(c echo.Context) error {
	var e Event
	c.Bind(&e)
	if e.Name == "" {
		return Error{HttpStatus: http.StatusBadRequest, Message: "Missing name of event."}
	}

	log.WithField("event", e).Debug("Adding new event")

	err := db.Events.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucket([]byte(e.Name))
		if err != nil {
			return Error{HttpStatus: http.StatusConflict, Message: "Event already exists", Data: e.Name, Internal: err}
		}
		return e.ToBolt(bucket)
	})

	if err != nil {
		return Error{
			HttpStatus: http.StatusInternalServerError,
			Message:    "Failed to create event",
			Data:       e.Name,
			Internal:   err,
		}
	}

	return c.JSON(http.StatusCreated, Error{http.StatusCreated, "Event created successfully", e, nil})
}

func GetRuns(c echo.Context) error {
	return db.Events.View(func(tx *bolt.Tx) error {
		runs := make([]Run, 0, 64) //Preallocate space for 64 runs.

		err := tx.ForEach(func(name []byte, b *bolt.Bucket) (err error) {
			var run Run
			err = run.FromBolt(b)
			runs = append(runs, run)
			return
		})
		if err != nil {

		}
		return c.JSON(http.StatusOK, runs)
	})
}

func PostRuns(c echo.Context) error {
	var r Run
	c.Bind(&r)

	if r.Game == "" || r.Event == "" || r.Category == "" {
		return Error{HttpStatus: http.StatusBadRequest, Message: "Game, Event and Category is required."}
	}

	err := db.Runs.Update(func(tx *bolt.Tx) error {
		for {
			r.Id = NewID()
			if tx.Bucket([]byte(r.Id)) == nil { //Check for existance
				break
			}
		}
		bucket, _ := tx.CreateBucket([]byte(r.Id))
		return r.ToBolt(bucket)

	})

	if err != nil {
		return Error{
			HttpStatus: http.StatusInternalServerError,
			Message:    "Failed to create run",
			Data:       r,
			Internal:   err,
		}
	}

	return nil
}
