package node

import (
	"encoding/json"
	"fmt"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"path/filepath"
	"sync"

	"github.com/tidwall/buntdb"
)

const hashName = "hash.db"

type hashCache struct {
	v   sync.Map
	db  *buntdb.DB
	cfg *config.Config
}

// DataHashInfo ...
type DataHashInfo struct {
	DataHash string               `json:"data_hash"`
	DataInfo core.MediaSerializer `json:"data_info"`
	AddrInfo core.AddrInfo        `json:"addr_info"`
}

func newDataHashInfo(data core.MediaSerializer) *DataHashInfo {
	return &DataHashInfo{
		DataHash: data.Hash(),
		DataInfo: data,
		AddrInfo: core.AddrInfo{},
	}
}

// Hash ...
func (v DataHashInfo) Hash() string {
	return v.DataHash
}

// Encode ...
func (v DataHashInfo) Encode() (string, error) {
	marshal, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}

// Decode ...
func (v *DataHashInfo) Decode(s string) error {
	err := json.Unmarshal([]byte(s), v)
	if err != nil {
		return fmt.Errorf("data decode failed:%w", err)
	}
	return nil
	//if v.dataInfo == nil {
	//	return fmt.Errorf("nil data info object")
	//}
	//err = v.dataInfo.Decode(v.DataInfo)
	//if err != nil {
	//	return err
	//}
	//if h := v.DataInfo.Hash(); h != v.DataHash {
	//	return fmt.Errorf("wrong hash(%s) from hash(%s)", h, v.DataHash)
	//}
	//return nil
}

func newHashCacher(cfg *config.Config) *hashCache {
	db, err := buntdb.Open(filepath.Join(cfg.Path, hashName))
	// Open the data.db file. It will be created if it doesn't exist.
	if err != nil {
		log.Fatal(err)
	}
	db.CreateIndex("hash", "*",
		buntdb.IndexJSON("data_info.root_hash"),
		buntdb.Desc(buntdb.IndexJSON("data_info.last_update")))
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
func (h *hashCache) Store(hash string, data core.DataEncoder) error {
	return h.db.Update(func(tx *buntdb.Tx) error {
		encode, err := data.Encode()
		if err != nil {
			return err
		}
		_, _, err = tx.Set(hash, encode, nil)
		return err
	})
}

// Load ...
func (h *hashCache) Load(hash string, data core.DataDecoder) error {
	return h.db.View(func(tx *buntdb.Tx) error {
		var datum string
		err := tx.Ascend("hash", func(key, value string) bool {
			fmt.Printf("%s: %s\n", key, value)
			datum = value
			return false
		})
		if err != nil {
			return err
		}
		return data.Decode(datum)
	})
}

// GC ...
func (h *hashCache) GC() error {
	if h.db != nil {
		if err := h.db.Shrink(); err != nil {
			return err
		}
	}
	return nil
}
