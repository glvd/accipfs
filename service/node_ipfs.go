package service

import (
	"github.com/glvd/accipfs/config"
	"github.com/ipfs/go-ipfs-http-client"
	"github.com/multiformats/go-multiaddr"
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
	ma, err := multiaddr.NewMultiaddr("path")
	if err != nil {
		return nil, err
	}
	api, err := httpapi.NewApi(ma)
	if err != nil {
		return nil, err
	}
	return &nodeIPFS{
		cfg: config,
		api: api,
	}, nil
}
