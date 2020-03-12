package service

import (
	"github.com/glvd/accipfs/config"
	"github.com/goextension/log"
	"os/exec"
	"path/filepath"
)

type nodeServerIPFS struct {
	cfg  *config.Config
	name string
	cmd  *exec.Cmd
}

// Start ...
func (n *nodeServerIPFS) Start() error {
	n.cmd = exec.Command(n.name, "daemon", "--routing", "none")
	err := n.cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

// Stop ...
func (n *nodeServerIPFS) Stop() error {
	return n.cmd.Process.Kill()
}

// Init ...
func (n *nodeServerIPFS) Init() error {
	cmd := exec.Command(n.name, "init", "--profile", "badgerds")
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command(n.name, "config", "Swarm.EnableAutoNATService", "--bool", "true")
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command(n.name, "config", "Swarm.EnableRelayHop", "--bool", "true")
	if err := cmd.Run(); err != nil {
		return err
	}
	log.Infow("ipfs init", "tag", outputHead)
	return nil
}

// NewNodeServerIPFS ...
func NewNodeServerIPFS(cfg config.Config) NodeServer {
	path := filepath.Join(cfg.Path, "bin", binName(cfg.ETH.Name))
	return &nodeServerIPFS{
		cfg:  &cfg,
		name: path,
	}
}
