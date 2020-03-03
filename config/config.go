package config

import (
	"encoding/json"
	"github.com/goextension/extmap"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
)

// IPFSConfig ...
type IPFSConfig struct {
	Name string
	Addr string
}

// ETHConfig ...
type ETHConfig struct {
	Name string
	Key  string
	Pass string
}

// AWSConfig ...
type AWSConfig struct {
	HostedZoneID       string
	RecordName         string
	AwsAccessKeyID     string
	AwsSecretAccessKey string
}

// Config ...
type Config struct {
	Path string
	ETH  ETHConfig
	IPFS IPFSConfig
	AWS  AWSConfig
}

var name = "config"
var ext = ".json"
var _config *Config

// Initialize ...
func Initialize() {
	cfg, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	_config = cfg
}

// LoadConfig ...
func LoadConfig() (*Config, error) {
	viper.AddConfigPath(currentPath())
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
	return ioutil.WriteFile(name+ext, by, 0755)
}

// Global ...
func Global() Config {
	return *_config
}

// Default ...
func Default() *Config {
	return &Config{
		Path: "data",
		ETH: ETHConfig{
			Name: "eth",
			Key:  "",
			Pass: "",
		},
		IPFS: IPFSConfig{
			Name: "ipfs",
			Addr: "",
		},
		AWS: AWSConfig{},
	}
}

func currentPath() string {
	dir, e := os.Getwd()
	if e != nil {
		return "."
	}
	return dir
}
