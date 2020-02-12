package config

type IPFSConfig struct {
	Addr string
}

type Config struct {
	IPFS               IPFSConfig
	AwsAccessKeyID     string
	AwsSecretAccessKey string
}
