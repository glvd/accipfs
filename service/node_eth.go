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

func (n *nodeClientETH) output(v ...interface{}) {
	fmt.Print(outputHead, " ")
	fmt.Print("[ETH]", " ")
	_, _ = n.out.Println(v...)
}

// Run ...
func (n *nodeClientETH) Run() {
	n.output("syncing node")
	if n.lock.Load() {
		n.output("node is already running")
		return
	}
	n.lock.Store(true)
	defer n.lock.Store(false)
	n.output("sync running")
	if !n.IsReady() {
		n.output("waiting for ready")
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
func (n *nodeClientETH) IsReady() bool {
	client, err := ethclient.Dial(n.cfg.ETH.Addr)
	if err != nil {
		log.Errorw("new node eth", "error", err)
		return false
	}
	n.client = client
	return true
}

// Node ...
func (n *nodeClientETH) Node() {

}

// Token ...
func (n *nodeClientETH) Token() {

}
