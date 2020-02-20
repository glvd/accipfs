package service

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/glvd/accipfs"
	"github.com/glvd/accipfs/config"
	"os/exec"
	"path/filepath"
)

const ethPath = ".ethereum"
const endPoint = "geth.ipc"

type nodeClientETH struct {
	cfg    config.Config
	client *ethclient.Client
}

func newETH(cfg config.Config) (*nodeClientETH, error) {
	client, e := ethclient.Dial(filepath.Join(cfg.Path, ethPath, endPoint))
	if e != nil {
		return nil, e
	}
	return &nodeClientETH{
		cfg:    cfg,
		client: client,
	}, nil
}

type nodeE struct {
	name string
	cmd  *exec.Cmd
}

// Start ...
func (n *nodeE) Start() {

}

// NodeServerETH ...
func NodeServerETH(cfg config.Config) Node {
	cmd := exec.Command(cfg.ETH.Name, "")
	cmd.Env = accipfs.Environ()
	return &nodeE{cmd: cmd}
}
