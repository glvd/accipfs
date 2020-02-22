package config

import (
	"github.com/goextension/extmap"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
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

var path = "config"
var ext = "json"

// LoadConfig ...
func LoadConfig() (*Config, error) {
	if path != "" {
		// Use config file from the flag.
		viper.SetConfigFile(path)
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
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("JSON")
	e := viper.MergeConfigMap(extmap.StructToMap(config))
	if e != nil {
		return e
	}
	viper.SetConfigFile(path)

	return viper.WriteConfig()
}
