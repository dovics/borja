package cache

import (
	"encoding/binary"
	"encoding/json"
	"time"

	bolt "go.etcd.io/bbolt"
)

type BoltDBCache struct {
	db *bolt.DB
}

func NewCache(path string) (*BoltDBCache, error) {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}

	return &BoltDBCache{
		db,
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

func (c *BoltDBCache) ClearBefore(t time.Time) error {
	timestamp := make([]byte, 8)
	binary.BigEndian.PutUint64(timestamp, uint64(t.Unix()))
	return c.db.Update(func(t *bolt.Tx) error {
		return t.ForEach(func(name []byte, b *bolt.Bucket) error {
			return b.ForEach(func(k, v []byte) error {
				if !CompareTimeStamp(k, timestamp) {
					return b.Delete(k)
				}

				return nil
			})

		})
	})
}

type TimeValue struct {
	Time  time.Time   `json:"time,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

func (c *BoltDBCache) GetAfter(t time.Time) (result []interface{}, err error) {

	err = c.db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			return b.ForEach(func(k, v []byte) error {
				timestamp := int64(binary.BigEndian.Uint64(k))
				if timestamp > t.Unix() {
					var value interface{}
					err := json.Unmarshal(v, &value)
					if err != nil {
						value = err
					}

					result = append(result, &TimeValue{
						Time:  time.Unix(timestamp, 0),
						Value: value,
					})
				}

				return nil
			})

		})
	})

	return
}

func CompareTimeStamp(a []byte, b []byte) bool {
	for i := 0; i < 8; i++ {
		if a[i] < b[i] {
			return true
		}
	}

	return false
}
