package node

import (
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"path/filepath"
	"sync"

	"github.com/tidwall/buntdb"
)

const hashName = "hashCache.db"

type hashCache struct {
	v   sync.Map
	db  *buntdb.DB
	cfg *config.Config
}

type hashData struct {
	Hash string
	Data core.DataInfoV1
}

func newHashCacher(cfg *config.Config) *hashCache {
	db, err := buntdb.Open(filepath.Join(cfg.Path, hashName))
	// Open the data.db file. It will be created if it doesn't exist.
	if err != nil {
		log.Fatal(err)
	}
	return &hashCache{
		cfg: cfg,
		db:  db,
	}
}

// Close ...
func (h *hashCache) Close() error {
	if h.db != nil {
		err := h.db.Close()
		if err != nil {
			return err
		}
		h.db = nil
		return nil
	}
	return nil
}

// Store ...
func (h *hashCache) Store(hash string, data core.DataInfoV1) {
	h.db.Update(func(tx *buntdb.Tx) error {
	})
}
