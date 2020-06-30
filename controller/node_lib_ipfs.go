package controller

import "github.com/glvd/accipfs/config"

type nodeLibIPFS struct {
}

func newNodeLibIPFS(cfg *config.Config) *nodeLibIPFS {
	return &nodeLibIPFS{}
}
