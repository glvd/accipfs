package service

import (
	"github.com/glvd/accipfs/config"
	"github.com/ipfs/go-ipfs-http-client"
)

type nodeIPFS struct {
	cfg config.Config
	api *httpapi.HttpApi
}

func newNodeIPFS(config config.Config) (*nodeIPFS, error) {
	return &nodeIPFS{
		cfg: config,
	}, nil
}
