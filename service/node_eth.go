package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/fatih/color"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/contract"
	"github.com/glvd/accipfs/contract/dtag"
	"github.com/glvd/accipfs/contract/node"
	"github.com/glvd/accipfs/contract/token"
	"github.com/glvd/accipfs/core"
	"github.com/goextension/log"
	"os"
	"sort"
	"strings"
	"time"
)

const ethPath = ".ethereum"
const endPoint = "geth.ipc"

type nodeETH struct {
	*serviceNode
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

// Run ...
func (n *nodeETH) Run() {
	if n.lock.Load() {
		output("service NodeClient is already running")
		return
	}
	n.lock.Store(true)
	defer n.lock.Store(false)
	if !n.IsReady() {
		output("waiting for ready")
		return
	}
	ctx := context.TODO()

	// get self node info
	nodeInfo, err := n.NodeInfo(ctx)
	if err != nil {
		log.Errorw("get eth node info", "tag", outputHead, "error", err.Error(), "node", nodeInfo)
		return
	}
	cnode := nodeInfo.Enode
	jsonString, _ := json.Marshal(nodeInfo.Protocols)
	var nodeProtocal ETHProtocol
	err = json.Unmarshal(jsonString, &nodeProtocal)
	if err != nil {
		return
	}
	ip := os.Getenv("IP")
	cnode = strings.Split(cnode, "@")[0] + "@" + ip + ":30303"
	//
	//// get active nodes
	var activePeers []string
	peers, err := n.AllPeers(ctx)
	if err != nil {
		output("get active eth node failed", err.Error())
	} else {
		output("get active eth nodes", len(peers))
	}
	for _, peer := range peers {
		jsStr, _ := json.Marshal(peer.Protocols)
		var peerProtocol ETHProtocol
		err := json.Unmarshal(jsStr, &peerProtocol)
		if err != nil {
			return
		}
		fmt.Println("peer difficulty", peerProtocol.Eth.Difficulty)
		// check if peers had enough blocks
		if float64(peerProtocol.Eth.Difficulty)/float64(nodeProtocal.Eth.Difficulty) > 0.9 {
			activePeers = append(activePeers, peer.Enode)
		}
	}

	// init contract
	cl := contract.Loader(n.cfg)

	// get decoded contract nodes
	err = cl.Node(func(node *node.AccelerateNode, opts *bind.TransactOpts) error {
		o := &bind.CallOpts{Pending: true}
		nodes, e := node.GetEthNodes(o)
		if e != nil {
			output("get contract node failed", e.Error())
			return e
		}
		output("get contract nodes", len(nodes))

		nodes = decodeNodes(n.cfg, nodes)
		//fmt.Println("[cPeers]", cPeers)
		// get decoded contract signer nodes
		masterNodes, e := node.GetSignerNodes(o)
		if e != nil {
			output("get contract node failed", e.Error())
			return e
		}
		output("get contract nodes", len(masterNodes))

		masterNodes = decodeNodes(n.cfg, masterNodes)
		// filter public network accessible nodes
		accessibleNodes := getAccessibleEthNodes(activePeers, "30303", 3*time.Second)
		// sync nodes
		newSignerNodes := difference([]string{cnode}, masterNodes)
		newAccNodes := difference(accessibleNodes, nodes)
		// node to be deleted
		deleteNodes := difference(nodes, getAccessibleEthNodes(nodes, "30303", 3*time.Second))
		var deleteIdx []int
		for _, dNode := range deleteNodes {
			for idx, cNode := range nodes {
				if cNode == dNode {
					deleteIdx = append(deleteIdx, idx)
				}
			}
		}

		// delete rest node
		if len(deleteIdx) > 0 {
			var err error
			sort.Sort(sort.Reverse(sort.IntSlice(deleteIdx)))
			for _, idx := range deleteIdx {
				_, err = node.DeleteEthNodes(opts, uint32(idx))
			}

			if err != nil {
				fmt.Println("<删除失效节点失败>", err.Error())
			} else {
				fmt.Println("[删除失效节点成功]")
			}
		}

		// crypto node info && add to contract
		if len(newAccNodes) > 0 {
			var err error
			fmt.Println("[adding node]", newAccNodes)
			for _, n := range encodeNodes(n.cfg, newAccNodes) {
				_, err = node.AddEthNodes(opts, []string{n})
			}
			if err != nil {
				fmt.Println("[add node failed]", err.Error())
			} else {
				fmt.Println("[add node success]")
			}
			// update gateway info
		} else {
			fmt.Println("[已经是最新节点数据]")
		}

		// add signer nodes
		if len(newSignerNodes) > 0 {
			fmt.Println("[adding signer node]", newSignerNodes)
			_, err := node.AddSignerNodes(opts, encodeNodes(n.cfg, newSignerNodes))
			if err != nil {
				fmt.Println("<添加主节点失败>", err.Error())
			} else {
				fmt.Println("[添加主节点成功]")
			}
		}
		vNodes := difference(accessibleNodes, masterNodes)
		mNodes := make(map[string]bool)
		for _, value := range vNodes {
			mNodes[value] = true
		}
		syncDNS(n.cfg, mNodes)
		return nil
	})

	if err != nil {
		log.Errorw("eth node process error", "tag", outputHead, "error", err)
		return
	}

	output("sync eth node complete")
	return
}

func newNodeETH(cfg *config.Config) (*nodeETH, error) {
	return &nodeETH{
		cfg:         cfg,
		serviceNode: nodeInstance(),
	}, nil
}

// IsReady ...
func (n *nodeETH) IsReady() bool {
	client, err := ethclient.Dial(config.ETHAddr())
	if err != nil {
		log.Errorw("new serviceNode eth", "tag", outputHead, "error", err)
		return false
	}
	n.client = client
	return true
}

// DMessage ...
func (n *nodeETH) DTag() (*dtag.DTag, error) {
	address := common.HexToAddress(n.cfg.ETH.MessageAddr)
	return dtag.NewDTag(address, n.client)
}

// NodeClient ...
func (n *nodeETH) Node() (*node.AccelerateNode, error) {
	address := common.HexToAddress(n.cfg.ETH.NodeAddr)
	return node.NewAccelerateNode(address, n.client)
}

// Token ...
func (n *nodeETH) Token() (*token.DhToken, error) {
	address := common.HexToAddress(n.cfg.ETH.TokenAddr)
	return token.NewDhToken(address, n.client)
}

// Peers ...
func (n *nodeETH) AllPeers(ctx context.Context) ([]ETHPeer, error) {
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
func (n *nodeETH) NodeInfo(ctx context.Context) (*core.ContractNode, error) {
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
func (n *nodeETH) AddPeer(ctx context.Context, peer string) error {
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
func (n *nodeETH) FindNo(ctx context.Context, no string) error {
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
