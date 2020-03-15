package cache

import (
	"github.com/glvd/accipfs/config"
	"github.com/gocacher/badger-cache/v2"
	"github.com/gocacher/cacher"
	"path/filepath"
)

type memoryCache struct {
	path   string
	memory map[string][]byte
	cache  cacher.Cacher
}

// Get ...
func (m memoryCache) Get(key string) ([]byte, error) {
	panic("implement me")
}

// GetD ...
func (m memoryCache) GetD(key string, v []byte) []byte {
	panic("implement me")
}

// Set ...
func (m memoryCache) Set(key string, val []byte) error {
	panic("implement me")
}

// SetWithTTL ...
func (m memoryCache) SetWithTTL(key string, val []byte, ttl int64) error {
	panic("implement me")
}

// Has ...
func (m memoryCache) Has(key string) (bool, error) {
	panic("implement me")
}

// Delete ...
func (m memoryCache) Delete(key string) error {
	panic("implement me")
}

// Clear ...
func (m memoryCache) Clear() error {
	panic("implement me")
}

// GetMultiple ...
func (m memoryCache) GetMultiple(keys ...string) (map[string][]byte, error) {
	panic("implement me")
}

// SetMultiple ...
func (m memoryCache) SetMultiple(values map[string][]byte) error {
	panic("implement me")
}

// DeleteMultiple ...
func (m memoryCache) DeleteMultiple(keys ...string) error {
	panic("implement me")
}

// New ...
func New(cfg config.Config) cacher.Cacher {
	cache.DefaultCachePath = filepath.Join(cfg.Path, ".cache")
	return &memoryCache{
		path:   cache.DefaultCachePath,
		memory: make(map[string][]byte),
		cache:  cache.New(),
	}

}
