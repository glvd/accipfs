package controller

import (
	"context"
	"fmt"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/glvd/accipfs/general"
	"github.com/goextension/io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

var _ core.ControllerService = &nodeBinETH{}

type nodeBinETH struct {
	ctx     context.Context
	cancel  context.CancelFunc
	cfg     *config.Config
	genesis *config.Genesis
	name    string
	cmd     *exec.Cmd
	msg     func(s string)
}

// MessageHandle ...
func (n *nodeBinETH) MessageHandle(f func(s string)) {
	if f != nil {
		n.msg = f
	}
}

// Stop ...
func (n *nodeBinETH) Stop() error {
	if n.cmd != nil {
		n.cancel()
		n.cmd = nil
	}
	return nil
}

// Start ...
func (n *nodeBinETH) Start() error {
	if core.NodeAccount.CompareInt(n.cfg.NodeType) {
		n.cmd = exec.CommandContext(n.ctx, n.name,
			"--datadir", config.DataDirETH(),
			"--networkid", strconv.FormatInt(n.genesis.Config.ChainID, 10),
			"--allow-insecure-unlock",
			"--rpccorsdomain", "*", "--rpc", "--rpcport", "8545", "--rpcaddr", "127.0.0.1",
			"--rpcapi", "admin,eth,net,web3,personal,miner",
			"--unlock", "54c0fa4a3d982656c51fe7dfbdcc21923a7678cb",
			"--password", filepath.Join(n.cfg.Path, "password"),
			"--mine", "--nodiscover",
		)
	} else {
		n.cmd = exec.CommandContext(n.ctx, n.name,
			"--datadir", config.DataDirETH(),
			"--networkid", strconv.FormatInt(n.genesis.Config.ChainID, 10),
			"--rpccorsdomain", "*", "--rpc", "--rpcport", "8545", "--rpcaddr", "127.0.0.1",
			"--rpcapi", "admin,eth,net,web3,personal,miner",
			"--mine", "--nodiscover",
		)
	}

	output("geth cmd: ", n.cmd.Args)
	pipe, err2 := n.cmd.StderrPipe()
	if err2 != nil {
		return err2
	}
	stdoutPipe, err2 := n.cmd.StdoutPipe()
	if err2 != nil {
		return err2
	}
	m := io.MultiReader(pipe, stdoutPipe)
	if n.cfg.ETH.LogOutput && n.msg != nil {
		go general.PipeReader(n.ctx, m, n.msg)
	}
	//else {
	//	go general.PipeDummy(n.ctx, module, m)
	//}

	err := n.cmd.Start()
	if err != nil {
		return err
	}
	//geth --datadir /root/.ethereum --miner.gasprice 1000 --targetgaslimit 50000000  --networkid 20190723 --allow-insecure-unlock --rpc --rpcaddr 0.0.0.0 --rpccorsdomain '*' --rpcapi db,eth,net,web3,personal --unlock 54C0fa4a3d982656c51fe7dFBdCc21923a7678cB --password /root/.ethereum/password --nodiscover --mine
	return nil
}

// Initialize ...
func (n *nodeBinETH) Initialize() error {
	_, err := os.Stat(config.DataDirETH())
	if err != nil && os.IsNotExist(err) {
		_ = os.MkdirAll(config.DataDirETH(), 0755)
	}
	cmd := exec.Command(n.name, "--datadir", config.DataDirETH(), "init", filepath.Join(n.cfg.Path, "genesis.json"))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("eth init:%w", err)
	}
	n.msg(string(out))
	return nil
}

func newNodeBinETH(cfg *config.Config) *nodeBinETH {
	path := filepath.Join(cfg.Path, "bin", general.BinName(cfg.ETH.Name))
	genesis, err := config.LoadGenesis(cfg)
	if err != nil {
		panic(err)
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &nodeBinETH{
		ctx:     ctx,
		cancel:  cancelFunc,
		cfg:     cfg,
		genesis: genesis,
		name:    path,
	}
}
