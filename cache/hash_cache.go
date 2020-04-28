package cache

import (
	"github.com/glvd/accipfs/config"
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
