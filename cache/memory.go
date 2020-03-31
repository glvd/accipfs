package cache

import (
	"container/ring"
	"encoding/json"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/gocacher/badger-cache/v2"
	"github.com/gocacher/cacher"
	"path/filepath"
	"sync"
)

// MemoryCache ...
type MemoryCache struct {
	path   string
	loop   ring.Ring
	memory map[string][]byte
	mut    sync.RWMutex
	cache  cacher.Cacher
}

func nodePrefix(name string) string {
	return "node_" + name
}

func hashPrefix(hash string) string {
	return "hash_" + hash
}

// Get ...
func (m *MemoryCache) Get(key string) ([]byte, error) {
	m.mut.RLock()
	defer m.mut.RUnlock()
	return m.cache.Get(key)
}

// GetD ...
func (m *MemoryCache) GetD(key string, v []byte) []byte {
	m.mut.RLock()
	defer m.mut.RUnlock()
	return m.cache.GetD(key, v)
}

// Set ...
func (m *MemoryCache) Set(key string, val []byte) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	return m.cache.Set(key, val)
}

// SetWithTTL ...
func (m *MemoryCache) SetWithTTL(key string, val []byte, ttl int64) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	return m.cache.SetWithTTL(key, val, ttl)
}

// Has ...
func (m *MemoryCache) Has(key string) (bool, error) {
	m.mut.RLock()
	defer m.mut.RUnlock()
	return m.cache.Has(key)
}

// Delete ...
func (m *MemoryCache) Delete(key string) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	return m.cache.Delete(key)
}

// Clear ...
func (m *MemoryCache) Clear() error {
	m.mut.Lock()
	defer m.mut.Unlock()
	return m.cache.Clear()
}

// GetMultiple ...
func (m *MemoryCache) GetMultiple(keys ...string) (map[string][]byte, error) {
	m.mut.RLock()
	defer m.mut.RUnlock()
	return m.cache.GetMultiple(keys...)
}

// SetMultiple ...
func (m *MemoryCache) SetMultiple(values map[string][]byte) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	return m.cache.SetMultiple(values)
}

// DeleteMultiple ...
func (m *MemoryCache) DeleteMultiple(keys ...string) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	return m.cache.DeleteMultiple(keys...)
}

// SetNodeInfo ...
func (m *MemoryCache) SetNodeInfo(info *core.NodeInfo) error {
	marshal, err := json.Marshal(info)
	if err != nil {
		return err
	}
	err = m.Set(nodePrefix(info.Name), marshal)
	if err != nil {
		return err
	}
	return nil
}

// GetNodeInfo ...
func (m *MemoryCache) GetNodeInfo(name string) (*core.NodeInfo, error) {
	var info core.NodeInfo
	get, err := m.Get(nodePrefix(name))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(get, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// New ...
func New(cfg *config.Config) *MemoryCache {
	cache.DefaultCachePath = filepath.Join(cfg.Path, ".cache")
	return &MemoryCache{
		path:   cache.DefaultCachePath,
		memory: make(map[string][]byte),
		cache:  cache.New(),
	}
}
