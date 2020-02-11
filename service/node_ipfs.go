package service

import (
	"github.com/glvd/accipfs/config"
	"github.com/ipfs/go-ipfs-http-client"
)

type nodeIPFS struct {
	cfg config.Config
	api *httpapi.HttpApi
}

// Start ...
func (n nodeIPFS) Start() {
	panic("implement me")
}

// NodeIPFS ...
func NodeIPFS(config config.Config) (Node, error) {
	return &nodeIPFS{
		cfg: config,
	}, nil
}
