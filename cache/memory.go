package cache

import (
	"container/ring"
	"github.com/glvd/accipfs/config"
	"github.com/gocacher/badger-cache/v2"
	"github.com/gocacher/cacher"
	"path/filepath"
	"sync"
)

type memoryCache struct {
	path   string
	loop   ring.Ring
	memory map[string][]byte
	mut    sync.RWMutex
	cache  cacher.Cacher
}

// Get ...
func (m *memoryCache) Get(key string) ([]byte, error) {
	m.mut.RLock()
	defer m.mut.RUnlock()
	return m.cache.Get(key)
}

// GetD ...
func (m *memoryCache) GetD(key string, v []byte) []byte {
	m.mut.RLock()
	defer m.mut.RUnlock()
	return m.cache.GetD(key, v)
}

// Set ...
func (m *memoryCache) Set(key string, val []byte) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	return m.cache.Set(key, val)
}

// SetWithTTL ...
func (m *memoryCache) SetWithTTL(key string, val []byte, ttl int64) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	return m.cache.SetWithTTL(key, val, ttl)
}

// Has ...
func (m *memoryCache) Has(key string) (bool, error) {
	m.mut.RLock()
	defer m.mut.RUnlock()
	return m.cache.Has(key)
}

// Delete ...
func (m *memoryCache) Delete(key string) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	return m.cache.Delete(key)
}

// Clear ...
func (m *memoryCache) Clear() error {
	m.mut.Lock()
	defer m.mut.Unlock()
	return m.cache.Clear()
}

// GetMultiple ...
func (m *memoryCache) GetMultiple(keys ...string) (map[string][]byte, error) {
	m.mut.RLock()
	defer m.mut.RUnlock()
	return m.cache.GetMultiple(keys...)
}

// SetMultiple ...
func (m *memoryCache) SetMultiple(values map[string][]byte) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	return m.cache.SetMultiple(values)
}

// DeleteMultiple ...
func (m *memoryCache) DeleteMultiple(keys ...string) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	return m.cache.DeleteMultiple(keys...)
}

// New ...
func New(cfg *config.Config) cacher.Cacher {
	cache.DefaultCachePath = filepath.Join(cfg.Path, ".cache")
	return &memoryCache{
		path:   cache.DefaultCachePath,
		memory: make(map[string][]byte),
		cache:  cache.New(),
	}
}
