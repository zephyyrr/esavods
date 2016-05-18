package models

import (
	"encoding/binary"
	"github.com/boltdb/bolt"
)

type Runs []Run

type Run struct {
	Id       Id                  `json:"id"`
	Game     string              `json:"game"`
	Players  StringList          `json:"players"`
	Category string              `json:"category"`
	Type     string              `json:"type"`
	Console  string              `json:"console"`
	Length   MilliSecondDuration `json:"length"`
	Event    string              `json:"event"`
	Tags     StringList          `json:"tags"`
	Vods     Vods                `json:"vods"`
}

func (r Run) ToBolt(bucket *bolt.Bucket) error {
	bucket.Put([]byte("Id"), []byte(r.Id))
	bucket.Put([]byte("Game"), []byte(r.Game))
	bucket.Put([]byte("Event"), []byte(r.Event))
	bucket.Put([]byte("Category"), []byte(r.Category))
	bucket.Put([]byte("Type"), []byte(r.Type))
	bucket.Put([]byte("Console"), []byte(r.Console))

	{ // Encode Length
		size := binary.Size(r.Length)
		buffer := make([]byte, size, size)
		binary.PutVarint(buffer, int64(r.Length))
		bucket.Put([]byte("Length"), buffer)
	}

	players, _ := bucket.CreateBucketIfNotExists([]byte("Players"))
	r.Players.ToBolt(players)

	tags, _ := bucket.CreateBucketIfNotExists([]byte("Tags"))
	r.Tags.ToBolt(tags)

	vods, _ := bucket.CreateBucketIfNotExists([]byte("Vods"))
	r.Vods.ToBolt(vods)

	return nil
}

func (r *Run) FromBolt(bucket *bolt.Bucket) error {
	r.Id = Id(bucket.Get([]byte("Id")))
	r.Game = string(bucket.Get([]byte("Game")))
	r.Category = string(bucket.Get([]byte("Category")))
	r.Type = string(bucket.Get([]byte("Type")))
	r.Console = string(bucket.Get([]byte("Console")))
	r.Event = string(bucket.Get([]byte("Event")))

	length, _ := binary.Varint(bucket.Get([]byte("Length")))
	r.Length = MilliSecondDuration(length)

	r.Players.FromBolt(bucket.Bucket([]byte("Players")))
	r.Tags.FromBolt(bucket.Bucket([]byte("Tags")))
	r.Vods.FromBolt(bucket.Bucket([]byte("Vods")))

	return nil
}

type StringList []string

func (pl StringList) ToBolt(bucket *bolt.Bucket) error {
	for i, player := range pl {
		bucket.Put([]byte{byte(i)}, []byte(player))
	}
	//bucket.Put([]byte("Length"), []byte{byte(len(pl))})
	return nil
}

func (pl *StringList) FromBolt(bucket *bolt.Bucket) error {
	/*
		size := bucket.Get([]byte("Length"))[0]
		list := make(PlayerList, 0, size)
		*pl = list

		for i := byte(0); i < size; i++ {
			player := string(bucket.Get([]byte{i}))
			*pl = append(*pl, player)
		}

		return nil
	*/

	bucket.ForEach(func(key, val []byte) error {
		*pl = append(*pl, string(val))
		return nil
	})
	return nil
}

type Vods map[Service]string

func (vods Vods) ToBolt(bucket *bolt.Bucket) error {
	for service, url := range vods {
		bucket.Put([]byte(service), []byte(url))
	}
	return nil
}

func (vods *Vods) FromBolt(bucket *bolt.Bucket) error {
	*vods = make(Vods)

	bucket.ForEach(func(key, val []byte) error {
		(*vods)[Service(key)] = string(val)
		return nil
	})

	return nil
}

type Service string
