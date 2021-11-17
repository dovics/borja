package cache

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"time"

	bolt "go.etcd.io/bbolt"
)

type BoltDBCache struct {
	path string
	db   *bolt.DB
}

func NewBlotCache(path string) (*BoltDBCache, error) {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}

	return &BoltDBCache{
		path: path,
		db:   db,
	}, nil
}

func (c *BoltDBCache) Add(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	now := make([]byte, 8)
	binary.BigEndian.PutUint64(now, uint64(time.Now().Unix()))
	c.db.Update(func(t *bolt.Tx) error {
		bucket, err := t.CreateBucketIfNotExists([]byte(key))
		if err != nil {
			return err
		}

		return bucket.Put(now, data)
	})
	return nil
}

func (c *BoltDBCache) Archive(name string) error {
	if err := c.db.Close(); err != nil {
		return err
	}

	err := os.Rename(c.path, name)
	if err != nil {
		return err
	}

	db, err := bolt.Open(c.path, 0666, nil)
	if err != nil {
		return err
	}

	c.db = db
	return nil
}
