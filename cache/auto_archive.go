package cache

import (
	"path/filepath"
	"time"

	"github.com/dovics/borja/util/log"
)

const timeFormat = "20060102-150405"

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
