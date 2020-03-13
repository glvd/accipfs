package service

import (
	"context"
	"fmt"
	"github.com/glvd/accipfs/config"
	"github.com/goextension/io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

type nodeServerETH struct {
	ctx     context.Context
	cancel  context.CancelFunc
	cfg     *config.Config
	genesis *config.Genesis
	name    string
	cmd     *exec.Cmd
}

// Stop ...
func (n *nodeServerETH) Stop() error {
	if n.cmd != nil {
		n.cancel()
		n.cmd = nil
	}
	return nil
}

// Start ...
func (n *nodeServerETH) Start() error {
	n.cmd = exec.CommandContext(n.ctx, n.name,
		"--datadir", config.DataDirETH(),
		"--networkid", strconv.FormatInt(n.genesis.Config.ChainID, 10),
		"--allow-insecure-unlock",
		"--rpccorsdomain", "*", "--rpc", "--rpcport", "8545", "--rpcaddr", "127.0.0.1",
		"--rpcapi", "admin,eth,net,web3,personal,miner",
		"--unlock", "945d35cd4a6549213e8d37feb5d708ec98906902",
		"--mine", "--nodiscover",
		"--password", filepath.Join(n.cfg.Path, "password"))
	fmt.Println("geth cmd: ", n.cmd.Args)
	pipe, err2 := n.cmd.StderrPipe()
	if err2 != nil {
		return err2
	}
	stdoutPipe, err2 := n.cmd.StdoutPipe()
	if err2 != nil {
		return err2
	}
	m := io.MultiReader(pipe, stdoutPipe)
	go screenOutput(n.ctx, m)
	err := n.cmd.Start()
	if err != nil {
		return err
	}
	//geth --datadir /root/.ethereum --miner.gasprice 1000 --targetgaslimit 50000000  --networkid 20190723 --allow-insecure-unlock --rpc --rpcaddr 0.0.0.0 --rpccorsdomain '*' --rpcapi db,eth,net,web3,personal --unlock 54C0fa4a3d982656c51fe7dFBdCc21923a7678cB --password /root/.ethereum/password --nodiscover --mine
	return nil
}

// Init ...
func (n *nodeServerETH) Init() error {
	_, err := os.Stat(config.DataDirETH())
	if err != nil && os.IsNotExist(err) {
		_ = os.MkdirAll(config.DataDirETH(), 0755)
	}
	cmd := exec.Command(n.name, "--datadir", config.DataDirETH(), "init", filepath.Join(n.cfg.Path, "genesis.json"))
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// NewNodeServerETH ...
func NewNodeServerETH(cfg config.Config) NodeServer {
	path := filepath.Join(cfg.Path, "bin", binName(cfg.ETH.Name))
	genesis, err := config.LoadGenesis(cfg)
	if err != nil {
		panic(err)
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &nodeServerETH{
		ctx:     ctx,
		cancel:  cancelFunc,
		cfg:     &cfg,
		genesis: genesis,
		name:    path,
	}
}
