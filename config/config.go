package config

import (
	"encoding/json"
	"github.com/mitchellh/go-homedir"
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
func LoadConfig() (*Config, error) {
	if DefaultPath != "" {
		// Use config file from the flag.
		viper.SetConfigFile(DefaultPath)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			return nil, err
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".json")
	}

	viper.AutomaticEnv()

	//if err := viper.ReadInConfig(); err == nil {
	//	fmt.Println("Using config file:", viper.ConfigFileUsed())
	//}

	cfg := Config{}

	err := viper.Unmarshal(&cfg)

	return &cfg, err
}

// SaveConfig ...
func SaveConfig(config *Config) error {
	by, e := json.Marshal(config)
	if e != nil {
		return e
	}
	if err := os.MkdirAll(DefaultPath, 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(DefaultPath, name+ext), by, 0755)
}
