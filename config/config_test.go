package config

import "testing"

func TestLoadConfig(t *testing.T) {
	t.Log(SaveConfig(&Config{
		Path:               "here",
		ETH:                ETHConfig{},
		IPFS:               IPFSConfig{},
		AwsAccessKeyID:     "",
		AwsSecretAccessKey: "",
	}))
	t.Log(LoadConfig())
}
