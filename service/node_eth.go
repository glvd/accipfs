package service

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
	"github.com/glvd/accipfs"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/contract/node"
	"github.com/glvd/accipfs/contract/token"
	"github.com/goextension/log"
	"os/exec"
)

const ethPath = ".ethereum"
const endPoint = "geth.ipc"

type nodeClientETH struct {
	*serviceNode
	cfg    config.Config
	client *ethclient.Client
	out    *color.Color
}

func (n *nodeClientETH) output(v ...interface{}) {
	v = append([]interface{}{outputHead, "[ETH]"}, v...)
	fmt.Println(v...)
}

// Run ...
func (n *nodeClientETH) Run() {
	n.output("syncing serviceNode")
	if n.lock.Load() {
		n.output("serviceNode is already running")
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
		cfg:         cfg,
		serviceNode: nodeInstance(),
		out:         color.New(color.FgRed),
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
		log.Errorw("new serviceNode eth", "error", err)
		return false
	}
	n.client = client
	return true
}

// Node ...
func (n *nodeClientETH) Node() (*node.AccelerateNode, error) {
	address := common.HexToAddress(n.cfg.ETH.NodeAddr)
	return node.NewAccelerateNode(address, n.client)
}

// Token ...
func (n *nodeClientETH) Token() (*token.DhToken, error) {
	address := common.HexToAddress(n.cfg.ETH.TokenAddr)
	return token.NewDhToken(address, n.client)
}
