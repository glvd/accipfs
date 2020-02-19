package service

import (
	"fmt"
	"github.com/glvd/accipfs"
	"github.com/glvd/accipfs/config"
	"os/exec"
)

// HandleInfo ...
type HandleInfo struct {
	ServiceName string
	Data        interface{}
	Callback    HandleCallback
}

// HandleCallback ...
type HandleCallback func(src interface{})

// Node ...
type Node interface {
	Start()
}

type nodeI struct {
	name string
	cmd  *exec.Cmd
}

type nodeE struct {
	name string
	cmd  *exec.Cmd
}

// Start ...
func (n *node) Start() {
	fmt.Println("starting", n.name)
}

// NodeI ...
func NodeI(cfg config.Config) Node {
	cmd := exec.Command(cfg.IPFS.Name, "")
	cmd.Env = accipfs.Environ()
	return &nodeI{cmd: cmd}
}

// NodeE ...
func NodeE(cfg config.Config) Node {
	cmd := exec.Command(cfg.ETH.Name, "")
	cmd.Env = accipfs.Environ()
	return &nodeE{cmd: cmd}
}
