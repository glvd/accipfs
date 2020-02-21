package config

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

// LoadConfig ...
func LoadConfig() *Config {
	return &Config{
		Path: "",
		ETH: ETHConfig{
			Name: "",
		},
		IPFS: IPFSConfig{
			Name: "",
			Addr: "",
		},
		AwsAccessKeyID:     "",
		AwsSecretAccessKey: "",
	}
}
