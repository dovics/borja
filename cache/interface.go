package cache

import "time"

type Cache interface {
	Add(key string, value interface{}) error
	ClearBefore(time.Time) error
	GetAfter(time.Time) ([]interface{}, error)
}
