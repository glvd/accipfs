package service

import (
	"context"
	"fmt"
	"github.com/glvd/accipfs/config"
	"github.com/goextension/log"
	"os"
	"os/exec"
	"path/filepath"
)

type nodeServerIPFS struct {
	ctx    context.Context
	cancel context.CancelFunc
	cfg    *config.Config
	name   string
	cmd    *exec.Cmd
}

// Start ...
func (n *nodeServerIPFS) Start() error {
	n.cmd = exec.CommandContext(n.ctx, n.name, "daemon", "--routing", "none")
	err := n.cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

// Stop ...
func (n *nodeServerIPFS) Stop() error {
	if n.cmd != nil {
		n.cancel()
		n.cmd = nil
	}
	return nil
}

// Init ...
func (n *nodeServerIPFS) Init() error {
	_, err := os.Stat(config.DataDirIPFS())
	if err != nil && os.IsNotExist(err) {
		_ = os.MkdirAll(config.DataDirIPFS(), 0755)
	}
	//os.Setenv("IPFS_PATH", filepath.Join(n.cfg.Path, config.DataDirIPFS()))
	cmd := exec.Command(n.name, "init", "--profile", "badgerds")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("init:%w", err)
	}
	log.Infow("ipfs init", "tag", outputHead, "log", string(out))

	cmd = exec.Command(n.name, "config", "Swarm.EnableAutoNATService", "--bool", "true")
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("config(nat):%w", err)
	}
	log.Infow("ipfs config set", "tag", outputHead, "log", string(out))

	cmd = exec.Command(n.name, "config", "Swarm.EnableRelayHop", "--bool", "true")
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("config(relay):%w", err)
	}
	log.Infow("ipfs init config set", "tag", outputHead, "log", string(out))
	log.Infow("ipfs init end", "tag", outputHead)
	return nil
}

// NewNodeServerIPFS ...
func NewNodeServerIPFS(cfg config.Config) NodeServer {
	path := filepath.Join(cfg.Path, "bin", binName(cfg.IPFS.Name))
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &nodeServerIPFS{
		ctx:    ctx,
		cancel: cancelFunc,
		cfg:    &cfg,
		name:   path,
	}
}
