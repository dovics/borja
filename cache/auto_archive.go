package cache

import (
	"os"
	"path/filepath"
	"time"

	"github.com/dovics/borja/util/log"
)

const timeFormat = "20060102-150405MST"

type AutoArchiveCache struct {
	Cache
	path        string
	autoArchive bool
	stop        chan struct{}
}

func AutoArchiveWrapper(c Cache, archivePath string, interval time.Duration) *AutoArchiveCache {
	cache := &AutoArchiveCache{Cache: c, path: archivePath}
	cache.StartAutoArchive(interval)
	return cache
}

func (c *AutoArchiveCache) IsAutoArchive() bool {
	return c.autoArchive
}

func (c *AutoArchiveCache) StartAutoArchive(interval time.Duration) {
	if c.autoArchive {
		c.StopAutoArchive()
	} else {
		c.autoArchive = true
	}

	c.stop = make(chan struct{})

	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-c.stop:
				return
			case t := <-ticker.C:
				if err := c.Archive(filepath.Join(c.path, t.Format(timeFormat))); err != nil {
					log.Error("file archive error: ", err)
				}
			}
		}
	}()
}

func (c *AutoArchiveCache) StopAutoArchive() {
	close(c.stop)
}

func (c *AutoArchiveCache) Clear(before time.Time) (interface{}, error) {
	entrys, err := os.ReadDir(c.path)
	for _, entry := range entrys {
		t, err := time.Parse(timeFormat, entry.Name())
		if err != nil {
			log.Error("file name time parse error: ", err)
			continue
		}

		if t.Before(before) {
			if err := os.Remove(filepath.Join(c.path, entry.Name())); err != nil {
				log.Error("file remove error: ", err)
				continue
			}
		}
	}

	return "Finish", err
}

func (c *AutoArchiveCache) ArchiveFiles() (interface{}, error) {
	entrys, err := os.ReadDir(c.path)
	var result []string
	for _, entry := range entrys {
		result = append(result, entry.Name())
	}

	return result, err
}
