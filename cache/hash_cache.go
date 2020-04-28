package cache

import (
	"encoding/json"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/gocacher/cacher"
	"sync"
)

// HashCache ...
type HashCache interface {
	Add(p string, s ...string) (e error)
	Get(p string) map[string]bool
	Has(p string) bool
}

type hashCache struct {
	prefix string
	mut    *sync.RWMutex
	cache  cacher.Cacher
	cfg    *config.Config
}

// Add ...
func (c *hashCache) Add(p string, s ...string) (e error) {
	p = prefixName(c.prefix, p)
	hashes := c.Get(p)
	changed := false
	for _, name := range s {
		if _, b := hashes[name]; b {
			continue
		}
		changed = true
		hashes[name] = true
	}
	if !changed {
		return nil
	}
	err := c.cache.Set(p, marshalMapNode(hashes))
	return err
}

// Get ...
func (c *hashCache) Get(p string) map[string]bool {
	p = prefixName(c.prefix, p)
	return unmarshalMapNode(c.cache.Get(p))
}

// Has ...
func (c *hashCache) Has(p string) bool {
	p = prefixName(c.prefix, p)
	b, e := c.cache.Has(p)
	if e != nil {
		logE("has check failed", "error", e, "bool", b)
		return false
	}
	return b
}

// NewHashCache ...
func NewHashCache(cfg *config.Config) HashCache {
	return &hashCache{
		cfg:    cfg,
		mut:    &sync.RWMutex{},
		prefix: "hash",
		cache:  New(cfg),
	}
}

// HashNodes ...
func (c *hashCache) HashNodes(hash string) []*core.Node {
	return nil
}

func marshalMapNode(node map[string]bool) []byte {
	marshal, err := json.Marshal(node)
	if err != nil {
		panic(err)
	}
	return marshal
}

func unmarshalMapNode(b []byte, err error) map[string]bool {
	node := map[string]bool{}
	if err != nil {
		logE("unmarshalMapNode data failed", "error", err)
		return node
	}
	err = json.Unmarshal(b, &node)
	if err != nil {
		logE("unmarshalMapNode failed", "error", err)
		return node
	}
	return node
}
