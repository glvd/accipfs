package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

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
