package node

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v2"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"path/filepath"
	"sync"
)

const hashName = "hashes"

type hashCache struct {
	v   sync.Map
	db  *badger.DB
	cfg *config.Config
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
	db, err := badger.Open(badger.DefaultOptions(filepath.Join(cfg.Path, hashName)))
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
