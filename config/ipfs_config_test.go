package config

import "testing"

func TestLoadIPFSServerConfig(t *testing.T) {
	WorkDir = `D:\workspace\pvt\data1\`
	cfg := Default()
	config, err := LoadIPFSServerConfig(cfg)
	if err != nil {
		return
	}
	t.Log(config)
}
