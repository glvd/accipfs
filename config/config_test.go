package config

import "testing"

func TestLoadConfig(t *testing.T) {
	t.Log(SaveConfig(&Config{
		Path: "test",
		ETH:  ETHConfig{},
		IPFS: IPFSConfig{},
	}))

	err := LoadConfig()
	if err != nil {
		panic(err)
	}

}
