package node

import (
	"encoding/json"
	"path/filepath"

	"github.com/dgraph-io/badger/v2"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
)

const (
	cacheDir = ".cache"
	hashName = "hashes"
	nodeName = "nodes"
)

type hashCache struct {
	//v   sync.Map
	db  *badger.DB
	cfg *config.Config
}
type nodeCache struct {
	db  *badger.DB
	cfg *config.Config
}

// Load ...
func (n *nodeCache) Load(hash string, data core.Unmarshaler) error {
	return n.db.View(
		func(txn *badger.Txn) error {
			item, err := txn.Get([]byte(hash))
			if err != nil {
				return err
			}
			return item.Value(func(val []byte) error {
				return data.Unmarshal(val)
			})
		})
}

// Store ...
func (n *nodeCache) Store(hash string, data core.Marshaler) error {
	return n.db.Update(
		func(txn *badger.Txn) error {
			encode, err := data.Marshal()
			if err != nil {
				return err
			}
			return txn.Set([]byte(hash), encode)
		})
}

// Close ...
func (n *nodeCache) Close() error {
	if n.db != nil {
		defer func() {
			n.db = nil
		}()
		return n.db.Close()
	}
	return nil
}

// Cacher ...
type Cacher interface {
	Load(hash string, data core.Unmarshaler) error
	Store(hash string, data core.Marshaler) error
	Close() error
}

// DataHashInfo ...
type DataHashInfo struct {
	DataHash string            `json:"data_hash"`
	DataInfo core.Serializable `json:"data_info"`
	AddrInfo core.AddrInfo     `json:"addr_info"`
}

func newDataHashInfo(data core.Serializable) *DataHashInfo {
	return &DataHashInfo{
		DataHash: data.Hash(),
		DataInfo: data,
	}
}

// Hash ...
func (v DataHashInfo) Hash() string {
	return v.DataHash
}

// Marshal ...
func (v DataHashInfo) Marshal() ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal ...
func (v *DataHashInfo) Unmarshal(b []byte) error {
	return json.Unmarshal(b, v)
}

func newHashCacher(cfg *config.Config) *hashCache {
	opts := badger.DefaultOptions(filepath.Join(cfg.Path, cacheDir, hashName))
	opts.Truncate = true
	db, err := badger.Open(opts)
	if err != nil {
		panic(err)
	}
	return &hashCache{
		cfg: cfg,
		db:  db,
	}
}

// Close ...
func (h *hashCache) Close() error {
	if h.db != nil {
		defer func() {
			h.db = nil
		}()
		return h.db.Close()
	}
	return nil
}

// Store ...
func (h *hashCache) Store(hash string, data core.Marshaler) error {
	return h.db.Update(
		func(txn *badger.Txn) error {
			encode, err := data.Marshal()
			if err != nil {
				return err
			}
			return txn.Set([]byte(hash), encode)
		})
}

// Load ...
func (h *hashCache) Load(hash string, data core.Unmarshaler) error {
	return h.db.View(
		func(txn *badger.Txn) error {
			item, err := txn.Get([]byte(hash))
			if err != nil {
				return err
			}
			return item.Value(func(val []byte) error {
				return data.Unmarshal(val)
			})
		})
}

// NewNodeCacher ...
func NewNodeCacher(cfg *config.Config) Cacher {
	opts := badger.DefaultOptions(filepath.Join(cfg.Path, cacheDir, nodeName))
	opts.Truncate = true
	db, err := badger.Open(opts)
	if err != nil {
		panic(err)
	}
	return &nodeCache{
		cfg: cfg,
		db:  db,
	}
}
