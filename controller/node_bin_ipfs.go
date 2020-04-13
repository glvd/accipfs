package controller

import (
	"context"
	"fmt"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/general"
	"github.com/goextension/io"
	"os"
	"os/exec"
	"path/filepath"
)

type nodeBinIPFS struct {
	ctx    context.Context
	cancel context.CancelFunc
	cfg    *config.Config
	name   string
	cmd    *exec.Cmd
}

// Start ...
func (n *nodeBinIPFS) Start() error {
	n.cmd = exec.CommandContext(n.ctx, n.name, "daemon", "--routing", "none")
	fmt.Println("ipfs cmd: ", n.cmd.Args)
	pipe, err2 := n.cmd.StderrPipe()
	if err2 != nil {
		return err2
	}
	stdoutPipe, err2 := n.cmd.StdoutPipe()
	if err2 != nil {
		return err2
	}
	m := io.MultiReader(pipe, stdoutPipe)
	go general.PipeScreen(n.ctx, module, m)
	err := n.cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

// Stop ...
func (n *nodeBinIPFS) Stop() error {
	if n.cmd != nil {
		n.cancel()
		n.cmd = nil
	}
	return nil
}

// Init ...
func (n *nodeBinIPFS) Init() error {
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
	logI("ipfs init", "log", string(out))
	cmd = exec.Command(n.name, "config", "Swarm.EnableAutoNATService", "--bool", "true")
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("config(nat):%w", err)
	}
	logI("ipfs init config set", "log", string(out))
	cmd = exec.Command(n.name, "config", "Swarm.EnableRelayHop", "--bool", "true")
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("config(relay):%w", err)
	}
	logI("ipfs init config set", "log", string(out))
	logI("ipfs init end")
	return nil
}

func newNodeBinIPFS(cfg *config.Config) *nodeBinIPFS {
	path := filepath.Join(cfg.Path, "bin", general.BinName(cfg.IPFS.Name))
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &nodeBinIPFS{
		ctx:    ctx,
		cancel: cancelFunc,
		cfg:    cfg,
		name:   path,
	}
}