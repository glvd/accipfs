package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const genesisData = `{
  "config": {
    "chainId": 20190723,
    "homesteadBlock": 1,
    "eip150Block": 2,
    "eip150Hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "eip155Block": 3,
    "eip158Block": 3,
    "byzantiumBlock": 4,
    "constantinopleBlock": 5,
    "clique": {
      "period": 10,
      "epoch": 30000
    }
  },
  "nonce": "0x0",
  "timestamp": "0x5d38141c",
  "extraData": "0x000000000000000000000000000000000000000000000000000000000000000054c0fa4a3d982656c51fe7dfbdcc21923a7678cb0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
  "gasLimit": "0xffffffff",
  "difficulty": "0x1",
  "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "coinbase": "0x0000000000000000000000000000000000000000",
  "alloc": {
    "20b98bEec5AE2e7149e70848a247406bE2c0cCA5": {
      "balance": "0x900000000000000000000000000000000000000000000000000000000000"
    },
    "2972Dd69A5242A4DF80f886b2fdD7a2DC99CD8A6": {
      "balance": "0x900000000000000000000000000000000000000000000000000000000000"
    },
    "945d35cd4a6549213e8D37Feb5d708EC98906902": {
      "balance": "0x900000000000000000000000000000000000000000000000000000000000"
    },
    "54c0fa4a3d982656c51fe7dfbdcc21923a7678cb": {
      "balance": "0x0"
    }
  },
  "number": "0x0",
  "gasUsed": "0x0",
  "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000"
}`

// Genesis ...
type Genesis struct {
	Config     GenesisConfig    `json:"config"`
	Nonce      string           `json:"nonce"`
	Timestamp  string           `json:"timestamp"`
	ExtraData  string           `json:"extraData"`
	GasLimit   string           `json:"gasLimit"`
	Difficulty string           `json:"difficulty"`
	MixHash    string           `json:"mixHash"`
	Coinbase   string           `json:"coinbase"`
	Alloc      map[string]Alloc `json:"alloc"`
	Number     string           `json:"number"`
	GasUsed    string           `json:"gasUsed"`
	ParentHash string           `json:"parentHash"`
}

// Alloc ...
type Alloc struct {
	Balance string `json:"balance"`
}

// GenesisConfig ...
type GenesisConfig struct {
	ChainID             int64  `json:"chainId"`
	HomesteadBlock      int64  `json:"homesteadBlock"`
	Eip150Block         int64  `json:"eip150Block"`
	Eip150Hash          string `json:"eip150Hash"`
	Eip155Block         int64  `json:"eip155Block"`
	Eip158Block         int64  `json:"eip158Block"`
	ByzantiumBlock      int64  `json:"byzantiumBlock"`
	ConstantinopleBlock int64  `json:"constantinopleBlock"`
	Clique              Clique `json:"clique"`
}

// Clique ...
type Clique struct {
	Period int64 `json:"period"`
	Epoch  int64 `json:"epoch"`
}

// LoadGenesis ...
func LoadGenesis(cfg *Config) (*Genesis, error) {
	var g Genesis
	path := filepath.Join(cfg.Path, "genesis.json")
	open, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(open)
	err = dec.Decode(&g)
	if err != nil {
		return nil, err
	}
	return &g, nil
}
