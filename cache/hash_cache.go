package cache

import (
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/gocacher/cacher"
)

// HashCache ...
type HashCache interface {
}

type hashCache struct {
	cache cacher.Cacher
}

// NewHashCache ...
func NewHashCache(cfg *config.Config) HashCache {
	return hashCache{
		cache: New(cfg),
	}
}

// AddHashed ...
func (c *hashCache) AddHashed() {

}

// HashNodes ...
func (c *hashCache) HashNodes(hash string) []*core.Node {
	return nil
}
