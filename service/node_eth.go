package service

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/glvd/accipfs/config"
	"path/filepath"
)

const ethPath = ".ethereum"
const endPoint = "geth.ipc"

type nodeETH struct {
	cfg    config.Config
	client *ethclient.Client
}

func newETH(cfg config.Config) (*nodeETH, error) {
	client, e := ethclient.Dial(filepath.Join(cfg.Path, ethPath, endPoint))
	if e != nil {
		return nil, e
	}
	return &nodeETH{
		cfg:    cfg,
		client: client,
	}, nil
}
