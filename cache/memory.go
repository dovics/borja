package cache

import (
	"time"
)

type memoryCache struct {
	m map[string]interface{}
}

var _ Cache = &memoryCache{
	m: make(map[string]interface{}),
}

func (c *memoryCache) Add(key string, value interface{}) error {
	c.m[key] = value
	return nil
}

func (c *memoryCache) ClearBefore(time.Time) error {
	c.m = make(map[string]interface{})
	return nil
}

func (c *memoryCache) GetAfter(time.Time) ([]interface{}, error) {
	return nil, nil
}
