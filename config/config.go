package config

// IPFSConfig ...
type IPFSConfig struct {
	Addr string
	Name string
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
