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
	Name    string
	Addr    string
	Timeout int
}

// ETHConfig ...
type ETHConfig struct {
	Name      string //bin name
	Addr      string //eth rpc address
	Key       string //eth key
	Pass      string //eth key pass
	NodeAddr  string //node contract address
	TokenAddr string //token contract address
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
	Path       string
	PrivateKey string
	PublicKey  string
	ETH        ETHConfig
	IPFS       IPFSConfig
	AWS        AWSConfig
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
			Name:      "geth",
			Addr:      DefaultGateway,
			Key:       "",
			Pass:      "",
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
