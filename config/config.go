package config

type IPFSConfig struct {
}

type Config struct {
	IPFS               IPFSConfig
	AwsAccessKeyID     string
	AwsSecretAccessKey string
}
