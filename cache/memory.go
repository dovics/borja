package cache

import (
	"encoding/json"
	"os"
)

type memoryCache struct {
	path string
	m    map[string]interface{}
}

func NewMemoryCache(basePath string) *memoryCache {
	return &memoryCache{
		path: basePath,
		m:    make(map[string]interface{}),
	}
}

func (c *memoryCache) Add(key string, value interface{}) error {
	c.m[key] = value
	return nil
}

func (c *memoryCache) Archive(name string) error {
	data, err := json.Marshal(c.m)
	if err != nil {
		return err
	}

	file, err := os.Open(name)
	if err != nil {
		return err
	}

	if _, err := file.Write(data); err != nil {
		return err
	}

	c.m = make(map[string]interface{})
	return nil
}
