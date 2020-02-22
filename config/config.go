package config

import (
	"encoding/json"
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

// DefaultPath ...
var DefaultPath = "config"
var name = "config"
var ext = ".json"

// LoadConfig ...
func LoadConfig() (interface{}, error) {
	// Use config file from the flag.
	viper.SetConfigFile(filepath.Join(DefaultPath, name+ext))
	//viper.AddConfigPath(DefaultConfigPath)
	//viper.SetConfigName(DefaultConfigName)

	err := viper.MergeInConfig()
	if err != nil {
		return nil, err
	}
	//m := extmap.ToMap(viper.AllSettings())
	//return m.Struct(&_config)
	viper.AutomaticEnv()

	var cfg interface{}

	err = viper.Unmarshal(&cfg)

	return &cfg, err
}

// SaveConfig ...
func SaveConfig(config *Config) error {
	by, e := json.MarshalIndent(config, "", " ")
	if e != nil {
		return e
	}
	if err := os.MkdirAll(DefaultPath, 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(DefaultPath, name+ext), by, 0755)
}
