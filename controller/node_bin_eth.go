package controller

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/contract/dtag"
	"github.com/glvd/accipfs/contract/node"
	"github.com/glvd/accipfs/contract/token"
	"github.com/glvd/accipfs/core"
	"github.com/glvd/accipfs/general"
	"github.com/goextension/io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var _ core.ControllerService = &nodeBinETH{}

// Network ...
type Network struct {
	Inbound       bool
	LocalAddress  string
	RemoteAddress string
	Static        bool
	Trusted       bool
}

// ETHPeer ...
type ETHPeer struct {
	Caps      []string
	ID        string
	Name      string
	Enode     string
	Network   Network
	Protocols interface{}
}

// Result ...
type Result struct {
	ID      string
	Jsonrpc string
	Result  []ETHPeer
}

// ETHProtocolInfo ...
type ETHProtocolInfo struct {
	Difficulty int    `json:"difficulty"`
	Head       string `json:"head"`
	Version    int    `json:"version"`
}

// ETHProtocol ...
type ETHProtocol struct {
	Eth ETHProtocolInfo `json:"eth"`
}
type nodeBinETH struct {
	ctx     context.Context
	cancel  context.CancelFunc
	cfg     *config.Config
	genesis *config.Genesis
	name    string
	cmd     *exec.Cmd
	msg     func(s string)
	client  *ethclient.Client
}

// MessageHandle ...
func (n *nodeBinETH) MessageHandle(f func(s string)) {
	if f != nil {
		n.msg = f
	}
}

// Msg ...
func (n *nodeBinETH) Msg(s string) {
	if n.msg != nil {
		n.msg(s)
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
	_, err := os.Stat(filepath.Join(n.cfg.Path, "password"))
	if core.NodeAccount.CompareInt(n.cfg.NodeType) && err == nil {
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

	err = n.cmd.Start()
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
	n.Msg(string(out))
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

// IsReady ...
func (n *nodeBinETH) IsReady() bool {
	client, err := ethclient.Dial(config.ETHAddr())
	if err != nil {
		log.Errorw("new serviceNode eth", "error", err)
		return false
	}
	n.client = client
	return true
}

// DMessage ...
func (n *nodeBinETH) DTag() (*dtag.DTag, error) {
	address := common.HexToAddress(n.cfg.ETH.DTagAddr)
	return dtag.NewDTag(address, n.client)
}

// NodeClient ...
func (n *nodeBinETH) Node() (*node.AccelerateNode, error) {
	address := common.HexToAddress(n.cfg.ETH.NodeAddr)
	return node.NewAccelerateNode(address, n.client)
}

// Token ...
func (n *nodeBinETH) Token() (*token.DhToken, error) {
	address := common.HexToAddress(n.cfg.ETH.TokenAddr)
	return token.NewDhToken(address, n.client)
}

// Peers ...
func (n *nodeBinETH) AllPeers(ctx context.Context) ([]ETHPeer, error) {
	var peers []ETHPeer
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	client, err := rpc.DialContext(cancelCtx, config.ETHAddr())
	if err != nil {
		return nil, err
	}
	defer client.Close()
	err = client.Call(&peers, "admin_peers")
	if err != nil {
		return nil, err
	}

	return peers, nil
}

// NewAccount ...
func (n *nodeBinETH) NodeInfo(ctx context.Context) (*core.ContractInfo, error) {
	var node core.ContractInfo
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	client, err := rpc.DialContext(cancelCtx, config.ETHAddr())
	if err != nil {
		return nil, err
	}
	defer client.Close()
	err = client.Call(&node, "admin_nodeInfo")
	if err != nil {
		return nil, err
	}

	return &node, nil
}

// AddPeer ...
func (n *nodeBinETH) AddPeer(ctx context.Context, peer string) error {
	var b bool
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	client, err := rpc.DialContext(cancelCtx, config.ETHAddr())
	if err != nil {
		return err
	}
	defer client.Close()
	err = client.Call(&b, "admin_addPeer", peer)
	if err != nil {
		return err
	}

	return nil
}

// FindNo ...
func (n *nodeBinETH) FindNo(ctx context.Context, no string) error {
	no = strings.ToUpper(no)
	t, err := n.DTag()
	if err != nil {
		return err
	}
	message, err := t.GetTagMessage(&bind.CallOpts{
		Pending: true,
		Context: ctx,
	}, "video", no)
	if err != nil {
		return err
	}
	fmt.Println("message", message.Value)

	return nil
}
