package models

import (
	"github.com/boltdb/bolt"
	"net/url"
	"time"
)

type Runs []Run

type Run struct {
	Id       Id            `json:"id"`
	Game     string        `json:"game"`
	Players  PlayerList    `json:"players"`
	Category string        `json:"category"`
	Type     string        `json:"type"`
	Console  string        `json:"console"`
	Length   time.Duration `json:"length"`
	Event    string        `json:"event"`
	Tags     []string      `json:"tags"`
	Vods     []Vod         `json:"vods"`
}

func (r Run) ToBolt(bucket *bolt.Bucket) error {
	bucket.Put([]byte("Id"), []byte(r.Id))
	bucket.Put([]byte("Game"), []byte(r.Game))
	bucket.Put([]byte("Event"), []byte(r.Event))
	bucket.Put([]byte("Category"), []byte(r.Category))
	bucket.Put([]byte("Type"), []byte(r.Type))
	bucket.Put([]byte("Console"), []byte(r.Console))
	bucket.Put([]byte("Length"), []byte(r.Length.String()))

	players, _ := bucket.CreateBucketIfNotExists([]byte("Players"))
	r.Players.ToBolt(players)

	return nil
}

func (r *Run) FromBolt(bucket *bolt.Bucket) error {
	r.Id = Id(bucket.Get([]byte("Id")))
	r.Game = string(bucket.Get([]byte("Game")))
	r.Category = string(bucket.Get([]byte("Category")))
	r.Type = string(bucket.Get([]byte("Type")))
	r.Console = string(bucket.Get([]byte("Console")))
	r.Event = string(bucket.Get([]byte("Event")))
	return r.Players.FromBolt(bucket.Bucket([]byte("Players")))
}

type PlayerList []string

func (pl PlayerList) ToBolt(bucket *bolt.Bucket) error {
	for i, player := range pl {
		bucket.Put([]byte{byte(i)}, []byte(player))
	}
	bucket.Put([]byte("Length"), []byte{byte(len(pl))})
	return nil
}

func (pl *PlayerList) FromBolt(bucket *bolt.Bucket) error {
	size := bucket.Get([]byte("Length"))[0]
	list := make(PlayerList, 0, size)
	*pl = list

	for i := byte(0); i < size; i++ {
		player := string(bucket.Get([]byte{i}))
		*pl = append(*pl, player)
	}

	return nil
}

type Vod struct {
	Service string  `json:"service"`
	URL     url.URL `json:"url"`
}
