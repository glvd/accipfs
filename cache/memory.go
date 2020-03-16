package cache

import (
	"container/ring"
	"github.com/glvd/accipfs/config"
	"github.com/gocacher/badger-cache/v2"
	"github.com/gocacher/cacher"
	"path/filepath"
)

type memoryCache struct {
	path   string
	loop   ring.Ring
	memory map[string][]byte
	cache  cacher.Cacher
}

// Get ...
func (m *memoryCache) Get(key string) ([]byte, error) {
	return m.cache.Get(key)
}

// GetD ...
func (m *memoryCache) GetD(key string, v []byte) []byte {
	return m.cache.GetD(key, v)
}

// Set ...
func (m *memoryCache) Set(key string, val []byte) error {
	return m.cache.Set(key, val)
}

// SetWithTTL ...
func (m *memoryCache) SetWithTTL(key string, val []byte, ttl int64) error {
	return m.cache.SetWithTTL(key, val, ttl)
}

// Has ...
func (m *memoryCache) Has(key string) (bool, error) {
	return m.cache.Has(key)
}

// Delete ...
func (m *memoryCache) Delete(key string) error {
	return m.cache.Delete(key)
}

// Clear ...
func (m *memoryCache) Clear() error {
	return m.cache.Clear()
}

// GetMultiple ...
func (m *memoryCache) GetMultiple(keys ...string) (map[string][]byte, error) {
	return m.cache.GetMultiple(keys...)
}

// SetMultiple ...
func (m *memoryCache) SetMultiple(values map[string][]byte) error {
	return m.cache.SetMultiple(values)
}

// DeleteMultiple ...
func (m *memoryCache) DeleteMultiple(keys ...string) error {
	return m.cache.DeleteMultiple(keys...)
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
