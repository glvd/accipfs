package service

import (
	"github.com/glvd/accipfs/config"
	"os"
	"os/exec"
	"path/filepath"
)

type nodeServerETH struct {
	cfg  *config.Config
	name string
	cmd  *exec.Cmd
}

// Start ...
func (n *nodeServerETH) Start() error {
}

// Init ...
func (n *nodeServerETH) Init() error {
	_, err := os.Stat(config.DataDirETH())
	if err != nil && os.IsNotExist(err) {
		_ = os.MkdirAll(config.DataDirETH(), 0755)
	}
	cmd := exec.Command(n.name, "--datadir", config.DataDirETH(), "init", filepath.Join(n.cfg.Path, "genesis.json"))
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// NewNodeServerETH ...
func NewNodeServerETH(cfg config.Config) NodeServer {
	path := filepath.Join(cfg.Path, "bin", binName(cfg.ETH.Name))
	return &nodeServerETH{
		cfg:  &cfg,
		name: path,
	}
}
