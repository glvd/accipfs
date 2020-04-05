package general

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// Alloc ...
type Alloc struct {
	Balance string `json:"balance"`
}

// Genesis ...
type Genesis struct {
	Config struct {
		ChainID             int `json:"chainId"`
		HomesteadBlock      int `json:"homesteadBlock"`
		Eip150Block         int `json:"eip150Block"`
		Eip155Block         int `json:"eip155Block"`
		Eip158Block         int `json:"eip158Block"`
		ByzantiumBlock      int `json:"byzantiumBlock"`
		ConstantinopleBlock int `json:"constantinopleBlock"`
		PetersburgBlock     int `json:"petersburgBlock"`
		Ethash              struct {
		} `json:"ethash"`
	} `json:"config"`
	Difficulty string            `json:"difficulty"`
	GasLimit   string            `json:"gasLimit"`
	Alloc      map[string]*Alloc `json:"alloc"`
}

// LoadGenesis ...
func LoadGenesis(closer io.ReadCloser) (*Genesis, error) {
	bytes, e := ioutil.ReadAll(closer)
	if e != nil {
		return nil, e
	}
	var g Genesis
	if err := json.Unmarshal(bytes, &g); err != nil {
		return nil, err
	}
	return &g, nil
}
