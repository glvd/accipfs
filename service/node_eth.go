package service

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
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

// Network ...
type Network struct {
	Inbound       bool
	LocalAddress  string
	RemoteAddress string
	Static        bool
	Trusted       bool
}

// Peer ...
type Peer struct {
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
	Result  []Peer
}

// Node self node info
type ETHNode struct {
	ID         string
	Enode      string
	IP         string
	Name       string
	ListenAddr string
	Ports      interface{}
	Protocols  interface{}
}

// NodeResult return node info
type ETHNodeResult struct {
	ID      string
	Jsonrpc string
	Result  Node
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

	// get self node info
	nodeInfo, err := eth.NodeInfo()
	if err != nil {
		fmt.Println("[获取本节点信息失败] ", err.Error())
	}
	//node := nodeInfo.Enode
	//jsonString, _ := json.Marshal(nodeInfo.Protocols)
	//var nodeProtocal EthProtocal
	//json.Unmarshal([]byte(jsonString), &nodeProtocal)
	//ip := os.Getenv("IP")
	//node = strings.Split(node, "@")[0] + "@" + ip + ":30303"
	//
	//// get active nodes
	//var activePeers []string
	//peers, err := eth.Peers()
	//if err != nil {
	//	fmt.Println("[获取活跃ETH节点失败] ", err.Error())
	//} else {
	//	fmt.Println("[当前活跃ETH节点数] ", len(activePeers))
	//}
	//for _, peer := range peers {
	//	jsStr, _ := json.Marshal(peer.Protocols)
	//	var peerProtocal EthProtocal
	//	json.Unmarshal([]byte(jsStr), &peerProtocal)
	//	fmt.Println("peer diffculty", peerProtocal.Eth.Difficulty)
	//	// check if peers had enough blocks
	//	if float64(peerProtocal.Eth.Difficulty)/float64(nodeProtocal.Eth.Difficulty) > 0.9 {
	//		activePeers = append(activePeers, peer.Enode)
	//	}
	//}
	//
	//// init contract
	//cl := eth.ContractLoader()
	//ac, auth, client := cl.AccelerateNode()
	//defer client.Close()
	//
	//// get decoded contract nodes
	//cPeers, err := ac.GetEthNodes(nil)
	//if err != nil {
	//	fmt.Println("[获取合约节点失败]", err.Error())
	//} else {
	//	fmt.Println("[合约已有节点数]", len(cPeers))
	//}
	//cPeers = decodeNodes(cPeers)
	//fmt.Println("[cPeers]", cPeers)
	//// get decoded contract signer nodes
	//cNodes, err := ac.GetSignerNodes(nil)
	//if err != nil {
	//	fmt.Println("[获取合约主节点失败] ", err.Error())
	//} else {
	//	fmt.Println("[合约已有主节点数]", len(cNodes))
	//}
	//cNodes = decodeNodes(cNodes)
	//fmt.Println("[cNodes]", cNodes)
	//// filter public network accessible nodes
	//accessibleNodes := getAccessibleEthNodes(activePeers, "30303")
	//// cDifference := difference(cPeers, accessibleNodes)
	//// sync nodes
	//newSignerNodes := difference([]string{node}, cNodes)
	//newAccNodes := difference(accessibleNodes, cPeers)
	//// node to be deleted
	//deleteNodes := difference(cPeers, getAccessibleEthNodes(cPeers, "30303"))
	//var deleteIdx []int
	//for _, dNode := range deleteNodes {
	//	for idx, cNode := range cPeers {
	//		if cNode == dNode {
	//			deleteIdx = append(deleteIdx, idx)
	//		}
	//	}
	//}
	//
	//// delete rest node
	//if len(deleteIdx) > 0 {
	//	var err error
	//	sort.Sort(sort.Reverse(sort.IntSlice(deleteIdx)))
	//	for _, idx := range deleteIdx {
	//		_, err = ac.DeleteEthNodes(auth, uint32(idx))
	//	}
	//
	//	if err != nil {
	//		fmt.Println("<删除失效节点失败>", err.Error())
	//	} else {
	//		fmt.Println("[删除失效节点成功]")
	//	}
	//}
	//
	//// crypto node info && add to contract
	//if len(newAccNodes) > 0 {
	//	var err error
	//	fmt.Println("[adding node]", newAccNodes)
	//	for _, n := range encodeNodes(newAccNodes) {
	//		_, err = ac.AddEthNodes(auth, []string{n})
	//	}
	//	if err != nil {
	//		fmt.Println("[add node failed]", err.Error())
	//	} else {
	//		fmt.Println("[add node success]")
	//	}
	//	// update gateway info
	//} else {
	//	fmt.Println("[已经是最新节点数据]")
	//}
	//
	//// add signer nodes
	//if len(newSignerNodes) > 0 {
	//	fmt.Println("[adding signer node]", newSignerNodes)
	//	_, err := ac.AddSignerNodes(auth, encodeNodes(newSignerNodes))
	//	if err != nil {
	//		fmt.Println("<添加主节点失败>", err.Error())
	//	} else {
	//		fmt.Println("[添加主节点成功]")
	//	}
	//}
	//
	//DNSSync(difference(accessibleNodes, cNodes))
	//
	//fmt.Println("<同步ETH节点完成>")
	return
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

// NodeInfo ...
func (n *nodeClientETH) ETHNodeInfo(ctx context.Context) (enode *ETHNode, e error) {
	var result ETHNode
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	cli, e := rpc.Dial(n.cfg.ETH.Addr)
	if e != nil {
		return nil, e
	}
	defer cli.Close()
	e = cli.Call(&result, "admin_nodeInfo")
	if e != nil {
		return nil, e
	}
	return &result, nil
}
