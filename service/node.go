package service

import (
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

type node struct {
	cmd *exec.Cmd
}

// Start ...
func (n *node) Start() {

}

// NodeI ...
func NodeI(cfg config.Config) Node {
	cmd := exec.Command(cfg.IPFS.Name, "")

	return &node{cmd: cmd}
}

// NodeE ...
func NodeE() {

}
