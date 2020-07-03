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
	db           *badger.DB
	iteratorOpts badger.IteratorOptions
	cfg          *config.Config
}

type nodeCache struct {
	db           *badger.DB
	iteratorOpts badger.IteratorOptions
	cfg          *config.Config
}

// Load ...
func (c *nodeCache) Load(hash string, data core.Unmarshaler) error {
	return c.db.View(
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
func (c *nodeCache) Store(hash string, data core.Marshaler) error {
	return c.db.Update(
		func(txn *badger.Txn) error {
			encode, err := data.Marshal()
			if err != nil {
				return err
			}
			return txn.Set([]byte(hash), encode)
		})
}

// Range ...
func (c *nodeCache) Range(f func(hash string, value string) bool) {
	c.db.View(func(txn *badger.Txn) error {
		iter := txn.NewIterator(c.iteratorOpts)
		defer iter.Close()
		var item *badger.Item
		var continueFlag bool
		for iter.Rewind(); iter.Valid(); iter.Next() {
			if !continueFlag {
				return nil
			}
			item = iter.Item()
			return item.Value(func(v []byte) error {
				key := item.Key()
				val, err := item.ValueCopy(v)
				if err != nil {
					return err
				}
				continueFlag = f(string(key), string(val))
				return nil
			})
		}
		return nil
	})
}

// Close ...
func (c *nodeCache) Close() error {
	if c.db != nil {
		defer func() {
			c.db = nil
		}()
		return c.db.Close()
	}
	return nil
}

// Cacher ...
type Cacher interface {
	Load(hash string, data core.Unmarshaler) error
	Store(hash string, data core.Marshaler) error
	Close() error
	Range(f func(hash string, value string) bool)
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

func hashCacher(cfg *config.Config) Cacher {
	opts := badger.DefaultOptions(filepath.Join(cfg.Path, cacheDir, hashName))
	opts.Truncate = true
	db, err := badger.Open(opts)
	if err != nil {
		panic(err)
	}
	itOpts := badger.DefaultIteratorOptions
	itOpts.Reverse = true
	return &hashCache{
		cfg:          cfg,
		iteratorOpts: itOpts,
		db:           db,
	}
}

// Store ...
func (c *hashCache) Store(hash string, data core.Marshaler) error {
	return c.db.Update(
		func(txn *badger.Txn) error {
			encode, err := data.Marshal()
			if err != nil {
				return err
			}
			return txn.Set([]byte(hash), encode)
		})
}

// Load ...
func (c *hashCache) Load(hash string, data core.Unmarshaler) error {
	return c.db.View(
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

// Range ...
func (c *hashCache) Range(f func(key, value string) bool) {
	c.db.View(func(txn *badger.Txn) error {
		iter := txn.NewIterator(c.iteratorOpts)
		defer iter.Close()
		var item *badger.Item
		var continueFlag bool
		for iter.Rewind(); iter.Valid(); iter.Next() {
			if !continueFlag {
				return nil
			}
			item = iter.Item()
			return item.Value(func(v []byte) error {
				key := item.Key()
				val, err := item.ValueCopy(v)
				if err != nil {
					return err
				}
				continueFlag = f(string(key), string(val))
				return nil
			})
		}
		return nil
	})
}

// Close ...
func (c *hashCache) Close() error {
	if c.db != nil {
		defer func() {
			c.db = nil
		}()
		return c.db.Close()
	}
	return nil
}

func nodeCacher(cfg *config.Config) Cacher {
	opts := badger.DefaultOptions(filepath.Join(cfg.Path, cacheDir, nodeName))
	opts.Truncate = true
	db, err := badger.Open(opts)
	if err != nil {
		panic(err)
	}
	itOpts := badger.DefaultIteratorOptions
	itOpts.Reverse = true
	return &nodeCache{
		cfg:          cfg,
		iteratorOpts: itOpts,
		db:           db,
	}
}
