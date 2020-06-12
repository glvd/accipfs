package config

import (
	"encoding/json"
	"fmt"
	"github.com/goextension/extmap"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const _keyDir = "key"
const _configName = "config"
const _configExt = ".json"
const _dataDirETH = ".eth"
const _dataDirIPFS = ".ipfs"
const _dataDirCache = ".cache"
const _localGateway = "http://127.0.0.1:%d"
const _ipfsGateway = "/ip4/127.0.0.1/tcp/%d"

// DefaultNodeContractAddr ...
var DefaultNodeContractAddr = "0xbaEEB7a3AF34a365ACAa1f8464A3374B58ac9889"

// DefaultTokenContractAddr ...
var DefaultTokenContractAddr = "0x9064322CfeE623A447ba5aF0dA6AD3341c073535"

// IPFSConfig ...
type IPFSConfig struct {
	LogOutput bool   `json:"log_output" mapstructure:"log_output"` //output log to screen
	Name      string `json:"name" mapstructure:"name"`
	//Addr    string `json:"addr" mapstructure:"addr"`
	Port    int `json:"port" mapstructure:"port"`
	Timeout int `json:"timeout" mapstructure:"timeout"`
}

// ETHConfig ...
type ETHConfig struct {
	LogOutput bool   `json:"log_output" mapstructure:"log_output"` //output log to screen
	Name      string `json:"name" mapstructure:"name"`             //bin name
	//Addr        string       `json:"addr" mapstructure:"addr"`                   //eth rpc address
	Port int `json:"port" mapstructure:"port"`
	//KeyHash     string                                    `json:"key_hash" mapstructure:"key_hash"`     //binary key hash
	NodeAddr    string `json:"node_addr" mapstructure:"node_addr"`       //node contract address
	TokenAddr   string `json:"token_addr" mapstructure:"token_addr"`     //token contract address
	MessageAddr string `json:"message_addr" mapstructure:"message_addr"` //dmessage contract address
	DTagAddr    string `json:"dtag_addr" mapstructure:"dtag_addr"`       //dtag contract address
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

// TLSCertificate ...
type TLSCertificate struct {
	KeyFile     string `json:"key_file" mapstructure:"key_file"`
	KeyPassFile string `json:"key_pass_file" mapstructure:"key_pass_file"`
}

// APIConfig ...
type APIConfig struct {
	Port    int            `json:"port" mapstructure:"port"`
	Version string         `json:"version" mapstructure:"version"`
	UseTLS  bool           `json:"use_tls" mapstructure:"use_tls"`
	TLS     TLSCertificate `json:"tls" mapstructure:"tls"`
	Timeout time.Duration  `json:"timeout" mapstructure:"timeout"`
}

// NodeConfig ...
type NodeConfig struct {
	Port          int           `json:"port" mapstructure:"port"`
	BindPort      int           `json:"bind_port" mapstructure:"bind_port"`
	BackupSeconds time.Duration `json:"backup_seconds" mapstructure:"backup_seconds"`
	PoolMax       int
}

// Config ...
type Config struct {
	Node       NodeConfig     `json:"node" mapstructure:"node"`
	API        APIConfig      `json:"api" mapstructure:"api"`
	UseTLS     bool           `json:"use_tls" mapstructure:"use_tls"`
	TLS        TLSCertificate `json:"tls" mapstructure:"tls"`
	Schema     string         `json:"schema" mapstructure:"schema"`
	Path       string         `json:"path" mapstructure:"path" `
	Account    string         `json:"account" mapstructure:"account"`
	PrivateKey string         `json:"private_key" mapstructure:"private_key"`
	PublicKey  string         `json:"public_key" mapstructure:"public_key"`
	ETH        ETHConfig      `json:"eth" mapstructure:"eth"`
	IPFS       IPFSConfig     `json:"ipfs" mapstructure:"ipfs"`
	AWS        AWSConfig      `json:"aws" mapstructure:"aws"`
	Interval   int64          `json:"interval" mapstructure:"interval"`
	NodeType   int            `json:"node_type" mapstructure:"node_type"`
	Limit      int64          `json:"limit" mapstructure:"limit"`
	Debug      bool           `json:"debug" mapstructure:"debug"`
	BootNode   []string       `json:"boot_node" mapstructure:"boot_node"`
}

// WorkDir ...
var WorkDir = ""

var _config *Config

func init() {
	WorkDir = currentPath()
}

// Initialize ...
func Initialize() {
	err := LoadConfig()
	if err != nil {
		panic(err)
	}
	err = os.Setenv("IPFS_PATH", DataDirIPFS())
	if err != nil {
		panic(err)
	}
}

// LoadConfig ...
func LoadConfig() error {
	viper.AddConfigPath(WorkDir)
	viper.SetConfigName(_configName)
	err := viper.MergeInConfig()
	if err != nil {
		return err
	}
	m := extmap.ToMap(viper.AllSettings())
	var cfg Config
	err = m.Struct(&cfg)
	if err != nil {
		return err
	}
	_config = &cfg
	fmt.Println("config loaded")
	return nil
}

// SaveConfig ...
func SaveConfig(config *Config) error {
	by, e := json.MarshalIndent(config, "", " ")
	if e != nil {
		return e
	}
	*_config = *config
	return ioutil.WriteFile(filepath.Join(WorkDir, _configName+_configExt), by, 0755)
}

// Global ...
func Global() Config {
	if _config == nil {
		panic("config must load first")
	}
	return *_config
}

// Default ...
func Default() *Config {
	def := &Config{
		Schema:     "http",
		Path:       WorkDir,
		Account:    "",
		PrivateKey: "",
		PublicKey:  "",
		ETH: ETHConfig{
			Name:      "geth",
			Port:      8545,
			NodeAddr:  DefaultNodeContractAddr,
			TokenAddr: DefaultTokenContractAddr,
		},
		IPFS: IPFSConfig{
			Name:    "ipfs",
			Port:    5001,
			Timeout: 30,
		},
		AWS:      AWSConfig{},
		Interval: 30,
		NodeType: 0x01,
		Limit:    500,
	}
	if _config == nil {
		_config = def
	}
	return def
}

// Initialize ...
func (c *Config) Init() error {
	err := os.Setenv("IPFS_PATH", filepath.Join(c.Path, _dataDirIPFS))
	if err != nil {
		return err
	}
	return nil
}

func (c Config) rpcAddr() string {
	return fmt.Sprintf("http://127.0.0.1:%d/rpc", c.Node.Port)
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
	return filepath.Join(Global().Path, _dataDirETH)
}

// KeyStoreDirETH ...
func KeyStoreDirETH() string {
	return filepath.Join(Global().Path, _dataDirETH, "keystore")
}

// DataDirIPFS ...
func DataDirIPFS() string {
	return filepath.Join(Global().Path, _dataDirIPFS)
}

// DataDirCache ...
func DataDirCache() string {
	return filepath.Join(Global().Path, _dataDirCache)
}

// KeyDir ...
func KeyDir() string {
	return filepath.Join(Global().Path, _keyDir)
}

// ETHAddr ...
func ETHAddr() string {
	return fmt.Sprintf(_localGateway, Global().ETH.Port)
}

// IPFSAddr ...
func IPFSAddr() string {
	return fmt.Sprintf(_ipfsGateway, Global().IPFS.Port)
}

// IPFSAddrHTTP ...
func IPFSAddrHTTP() string {
	return fmt.Sprintf(_localGateway, Global().IPFS.Port)
}

// RPCAddr ...
func RPCAddr() string {
	return Global().rpcAddr()
}
