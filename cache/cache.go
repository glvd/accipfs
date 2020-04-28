package cache

import (
	"github.com/glvd/accipfs/config"
	cache "github.com/gocacher/badger-cache/v2"
	"github.com/gocacher/cacher"
	"path/filepath"
)

type _cache struct {
	cacher.Cacher
	path string
}

// New ...
func New(cfg *config.Config) cacher.Cacher {
	cache.DefaultCachePath = filepath.Join(cfg.Path, ".cache")
	return &_cache{
		path:   cache.DefaultCachePath,
		Cacher: cache.New(),
	}
}
