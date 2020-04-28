package cache

import (
	"fmt"
	"github.com/glvd/accipfs/config"
	cache "github.com/gocacher/badger-cache/v2"
	"github.com/gocacher/cacher"
	"path/filepath"
)

type _cache struct {
	cacher.Cacher
	path string
}

var _instance *_cache

// InitCache ...
func InitCache(cfg *config.Config) {
	if _instance == nil {
		cache.DefaultCachePath = filepath.Join(cfg.Path, ".cache")
		//cacher.Register(cache.New())
		_instance = &_cache{
			path:   cache.DefaultCachePath,
			Cacher: cache.New(),
		}
	}

}

// New ...
func New(cfg *config.Config) cacher.Cacher {
	if _instance == nil {
		InitCache(cfg)
	}
	return _instance
}

func prefixName(p string, n string) string {
	return fmt.Sprintf("%s_%s", p, n)
}
