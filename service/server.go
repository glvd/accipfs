package service

import (
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/general"
	"os/exec"
	"path/filepath"
)

// NodeServer ...
type NodeServer interface {
	Start() error
	Init() error
}

// Server ...
type Server struct {
	cfg *config.Config
}

// NewServer ...
func NewServer(cfg config.Config) *Server {
	return &Server{cfg: &cfg}
}

type nodeServerETH struct {
	cfg  *config.Config
	name string
	cmd  *exec.Cmd
}

// Start ...
func (n *nodeServerETH) Start() {
	panic("TODO")
}

// Init ...
func (n *nodeServerETH) Init() {
	exec.Command(n.name, "--datadir", "../.geth", "init", "../contract")
}

// NewNodeServerETH ...
func NewNodeServerETH(cfg config.Config) NodeClient {
	path := filepath.Join(general.CurrentDir(), "bin", cfg.ETH.Name)
	return &nodeServerETH{
		name: path,
	}
}
