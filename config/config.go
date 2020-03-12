package config

import (
	"encoding/json"
	"github.com/goextension/extmap"
	"github.com/goextension/log"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
)

// IPFSConfig ...
type IPFSConfig struct {
	Name    string `json:"name" mapstructure:"name"`
	Addr    string `json:"addr" mapstructure:"addr"`
	Timeout int    `json:"timeout" mapstructure:"timeout"`
}

// ETHConfig ...
type ETHConfig struct {
	ETHKeyFile  `json:"key_file" mapstructure:"key_file"` //default key file
	KeyFileList []ETHKeyFile                              `json:"key_file_list" mapstructure:"key_file_list"` //key file list
	Name        string                                    `json:"name" mapstructure:"name"`                   //bin name
	Addr        string                                    `json:"addr" mapstructure:"addr"`                   //eth rpc address
	KeyHash     string                                    `json:"key_hash" mapstructure:"key_hash"`           //binary key hash
	NodeAddr    string                                    `json:"node_addr" mapstructure:"node_addr"`         //node contract address
	TokenAddr   string                                    `json:"token_addr" mapstructure:"token_addr"`       //token contract address
}

// AWSConfig ...
type AWSConfig struct {
	HostedZoneID       string `json:"hosted_zone_id" mapstructure:"hosted_zone_id"`
	RecordName         string `json:"record_name" mapstructure:"record_name"`
	AwsAccessKeyID     string `json:"aws_access_key_id" mapstructure:"aws_access_key_id"`
	AwsSecretAccessKey string `json:"aws_secret_access_key" mapstructure:"aws_secret_access_key"`
}

// ETHKeyFile ...
type ETHKeyFile struct {
	Name string `json:"name" mapstructure:"name"`
	Pass string `json:"pass" mapstructure:"pass"`
}

// Config ...
type Config struct {
	Path       string     `json:"path" mapstructure:"path" `
	PrivateKey string     `json:"private_key" mapstructure:"private_key"`
	PublicKey  string     `json:"public_key" mapstructure:"public_key"`
	ETH        ETHConfig  `json:"eth" mapstructure:"eth"`
	IPFS       IPFSConfig `json:"ipfs" mapstructure:"ipfs"`
	AWS        AWSConfig  `json:"aws" mapstructure:"aws"`
}

var name = "config"
var ext = ".json"
var _config *Config

// DefaultGateway ...
var DefaultGateway = "http://127.0.0.1:8545"

// DefaultNodeContractAddr ...
var DefaultNodeContractAddr = "0xbaEEB7a3AF34a365ACAa1f8464A3374B58ac9889"

// DefaultTokenContractAddr ...
var DefaultTokenContractAddr = "0x9064322CfeE623A447ba5aF0dA6AD3341c073535"

// WorkDir ...
var WorkDir = ""

var dataDirETH = ".eth"

var dataDirIPFS = ".ipfs"

func init() {
	WorkDir = currentPath()
}

// Initialize ...
func Initialize() {
	cfg, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	_config = cfg
	err = os.Setenv("IPFS_PATH", DataDirIPFS())
	if err != nil {
		panic(err)
	}
}

// LoadConfig ...
func LoadConfig() (*Config, error) {
	viper.AddConfigPath(WorkDir)
	viper.SetConfigName(name)
	err := viper.MergeInConfig()
	if err != nil {
		return nil, err
	}
	m := extmap.ToMap(viper.AllSettings())

	var cfg Config
	log.Infof("cfg map:%+v", m)
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
	return ioutil.WriteFile(filepath.Join(WorkDir, name+ext), by, 0755)
}

// Global ...
func Global() Config {
	return *_config
}

// Default ...
func Default() *Config {
	return &Config{
		Path: WorkDir,
		ETH: ETHConfig{
			Name:    "geth",
			Addr:    DefaultGateway,
			KeyHash: "",
			ETHKeyFile: ETHKeyFile{
				Name: "",
				Pass: "",
			},
			NodeAddr:  DefaultNodeContractAddr,
			TokenAddr: DefaultTokenContractAddr,
		},
		IPFS: IPFSConfig{
			Name:    "ipfs",
			Addr:    "/ip4/127.0.0.1/tcp/5001",
			Timeout: 30,
		},
	}
}

func currentPath() string {
	dir, e := os.Getwd()
	if e != nil {
		return "."
	}
	return dir
}

// DataDirETH ...
func DataDirETH() string {
	return filepath.Join(Global().Path, dataDirETH)
}

// DataDirIPFS ...
func DataDirIPFS() string {
	return filepath.Join(Global().Path, dataDirIPFS)
}
