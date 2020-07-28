package controller

import (
	"context"
	"github.com/glvd/accipfs/config"
)

type nodeLibIPFS struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func newNodeLibIPFS(cfg *config.Config) *nodeLibIPFS {
	ctx, cancel := context.WithCancel(context.Background())
	return &nodeLibIPFS{
		ctx:    ctx,
		cancel: cancel,
	}
}
