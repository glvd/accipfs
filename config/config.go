package config

import (
	"encoding/json"
	"github.com/goextension/extmap"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
)

// IPFSConfig ...
type IPFSConfig struct {
	Name string
	Addr string
}

// ETHConfig ...
type ETHConfig struct {
	Name string
}

// Config ...
type Config struct {
	Path               string
	ETH                ETHConfig
	IPFS               IPFSConfig
	AwsAccessKeyID     string
	AwsSecretAccessKey string
}

// DefaultConfigPath ...
var DefaultConfigPath = "config"
var name = "config"
var ext = ".json"

// LoadConfig ...
func LoadConfig() (*Config, error) {
	viper.AddConfigPath(filepath.Join(DefaultConfigPath))
	viper.SetConfigName(name)

	err := viper.MergeInConfig()
	if err != nil {
		return nil, err
	}
	m := extmap.ToMap(viper.AllSettings())

	var cfg Config
	err = m.Struct(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SaveConfig ...
func SaveConfig(config *Config) error {
	by, e := json.MarshalIndent(config, "", " ")
	if e != nil {
		return e
	}

	if err := os.MkdirAll(DefaultConfigPath, 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(DefaultConfigPath, name+ext), by, 0755)
}
