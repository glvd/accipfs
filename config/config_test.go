package config

import "testing"

func TestLoadConfig(t *testing.T) {
	t.Log(SaveConfig(&Config{
		Path:               "test",
		ETH:                ETHConfig{},
		IPFS:               IPFSConfig{},
		AwsAccessKeyID:     "",
		AwsSecretAccessKey: "",
	}))

	v, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	t.Logf("%+v", v)
}
