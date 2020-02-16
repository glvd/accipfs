package service

import (
	"github.com/glvd/accipfs/config"
)

const ethPath = ".ethereum"
const endPoint = "geth.ipc"

type nodeETH struct {
	cfg config.Config
}

func newETH(cfg config.Config) *nodeETH {
	return &nodeETH{cfg: cfg}
}
