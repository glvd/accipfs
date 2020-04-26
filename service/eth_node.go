package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/contract/dtag"
	"github.com/glvd/accipfs/contract/node"
	"github.com/glvd/accipfs/contract/token"
	"github.com/glvd/accipfs/core"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/fatih/color"
)

const ethPath = ".ethereum"
const ethEndPoint = "geth.ipc"

type ethNode struct {
	cfg    *config.Config
	client *ethclient.Client
	out    *color.Color
}

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

func newNodeETH(cfg *config.Config) (*ethNode, error) {
	return &ethNode{
		cfg: cfg,
	}, nil
}

// IsReady ...
func (n *ethNode) IsReady() bool {
	client, err := ethclient.Dial(config.ETHAddr())
	if err != nil {
		logE("new serviceNode eth", "error", err)
		return false
	}
	n.client = client
	return true
}

// DMessage ...
func (n *ethNode) DTag() (*dtag.DTag, error) {
	address := common.HexToAddress(n.cfg.ETH.DTagAddr)
	return dtag.NewDTag(address, n.client)
}

// NodeClient ...
func (n *ethNode) Node() (*node.AccelerateNode, error) {
	address := common.HexToAddress(n.cfg.ETH.NodeAddr)
	return node.NewAccelerateNode(address, n.client)
}

// Token ...
func (n *ethNode) Token() (*token.DhToken, error) {
	address := common.HexToAddress(n.cfg.ETH.TokenAddr)
	return token.NewDhToken(address, n.client)
}

// Peers ...
func (n *ethNode) AllPeers(ctx context.Context) ([]ETHPeer, error) {
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
func (n *ethNode) NodeInfo(ctx context.Context) (*core.ContractNode, error) {
	var node core.ContractNode
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
func (n *ethNode) AddPeer(ctx context.Context, peer string) error {
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
func (n *ethNode) FindNo(ctx context.Context, no string) error {
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
