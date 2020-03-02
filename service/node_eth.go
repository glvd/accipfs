package service

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/glvd/accipfs"
	"github.com/glvd/accipfs/config"
	"os/exec"
)

const ethPath = ".ethereum"
const endPoint = "geth.ipc"

type nodeClientETH struct {
	cfg    config.Config
	client *ethclient.Client
}

func (e *nodeClientETH) output(v ...interface{}) {
	fmt.Println(append([]interface{}{outputHead, "[ETH]"}, v...)...)
}

// Run ...
func (e *nodeClientETH) Run() {
	e.output("eth running")
}

func newETH(cfg config.Config) (*nodeClientETH, error) {
	//client, e := ethclient.Dial(filepath.Join(cfg.Path, ethPath, endPoint))
	//if e != nil {
	//	return nil, e
	//}
	return &nodeClientETH{
		cfg: cfg,
		//client: client,
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
	//TODO
	return false
}

// Node ...
func (e *nodeClientETH) Node() {

}

// Token ...
func (e *nodeClientETH) Token() {

}
