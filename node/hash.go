package node

import (
	"github.com/glvd/accipfs/config"
	"path/filepath"
	"sync"

	"github.com/tidwall/buntdb"
)

const hashName = "hash.db"

type hash struct {
	path string
	v    sync.Map
	db   *buntdb.DB
}

func newHashCacher(cfg *config.Config) *hash {
	db, err := buntdb.Open(filepath.Join(cfg.Path, hashName))
	// Open the data.db file. It will be created if it doesn't exist.
	if err != nil {
		log.Fatal(err)
	}
	return &hash{
		db: db,
	}
}

// Close ...
func (h *hash) Close() error {
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
