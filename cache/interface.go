package cache

type Cache interface {
	Add(key string, value interface{}) error
	Archive(name string) error
}
