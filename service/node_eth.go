package service

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
	"github.com/glvd/accipfs"
	"github.com/glvd/accipfs/config"
	"github.com/goextension/log"
	"os/exec"
)

const ethPath = ".ethereum"
const endPoint = "geth.ipc"

type nodeClientETH struct {
	*node
	cfg    config.Config
	client *ethclient.Client
	out    *color.Color
}

func (e *nodeClientETH) output(v ...interface{}) {
	fmt.Print(outputHead, " ")
	fmt.Print("[ETH]", " ")
	_, _ = e.out.Println(v...)
}

// Run ...
func (e *nodeClientETH) Run() {
	e.output("syncing node")
	if e.lock.Load() {
		e.output("ipfs node is already running")
		return
	}
	e.lock.Store(true)
	defer e.lock.Store(false)
	e.output("ipfs sync running")
	if !e.IsReady() {
		e.output("waiting for eth ready")
		return
	}
}

func newETH(cfg config.Config) (*nodeClientETH, error) {

	return &nodeClientETH{
		cfg:  cfg,
		node: nodeInstance(),
		out:  color.New(color.FgRed),
	}, nil
}

type nodeServerETH struct {
	name string
	cmd  *exec.Cmd
}

// Start ...
func (n *nodeServerETH) Start() {
	panic("TODO")
}

// NodeServerETH ...
func NodeServerETH(cfg config.Config) Node {
	cmd := exec.Command(cfg.ETH.Name, "")
	cmd.Env = accipfs.Environ()
	return &nodeServerETH{cmd: cmd}
}

// IsReady ...
func (e *nodeClientETH) IsReady() bool {
	client, err := ethclient.Dial(e.cfg.ETH.Addr)
	if err != nil {
		log.Errorw("new node eth", "error", err)
		return false
	}
	e.client = client
	return true
}

// Node ...
func (e *nodeClientETH) Node() {

}

// Token ...
func (e *nodeClientETH) Token() {

}
