package config

// IPFSConfig ...
type IPFSConfig struct {
	Addr string
	Name string
}

// Config ...
type Config struct {
	Path               string
	IPFS               IPFSConfig
	AwsAccessKeyID     string
	AwsSecretAccessKey string
}
