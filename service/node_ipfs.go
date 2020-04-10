package service

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/glvd/accipfs/contract/node"
	"github.com/glvd/accipfs/core"
	"github.com/goextension/log"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"net"
	"sort"
	"strings"
	"time"

	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/contract"
	"github.com/ipfs/go-ipfs-http-client"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
)

const ipfsPath = ".ipfs"
const ipfsAPI = "api"

type nodeIPFS struct {
	*serviceNode
	cfg *config.Config
	api *httpapi.HttpApi
}

// PeerID ...
type PeerID struct {
	Addresses       []string `json:"Addresses"`
	AgentVersion    string   `json:"AgentVersion"`
	ID              string   `json:"ID"`
	ProtocolVersion string   `json:"ProtocolVersion"`
	PublicKey       string   `json:"PublicKey"`
}

// NewNodeIPFS ...
func NewNodeIPFS(cfg *config.Config) (Node, error) {
	return newNodeIPFS(cfg)
}

func newNodeIPFS(cfg *config.Config) (*nodeIPFS, error) {
	node := &nodeIPFS{
		cfg:         cfg,
		serviceNode: nodeInstance(),
	}
	if err := node.connect(); err != nil {
		return nil, err
	}
	return node, nil
}

// SwarmConnect ...
func (n *nodeIPFS) SwarmConnect(ctx context.Context, addr string) (e error) {
	ma, e := multiaddr.NewMultiaddr(addr)
	if e != nil {
		return e
	}
	info, e := peer.AddrInfoFromP2pAddr(ma)
	if e != nil {
		return e
	}
	e = n.api.Swarm().Connect(ctx, *info)
	if e != nil {
		return e
	}
	return nil
}

func (n *nodeIPFS) connect() (e error) {
	ma, err := multiaddr.NewMultiaddr(config.IPFSAddr())
	if err != nil {
		return err
	}
	n.api, e = httpapi.NewApi(ma)
	return
}

// SwarmPeers ...
func (n *nodeIPFS) SwarmPeers(ctx context.Context) ([]iface.ConnectionInfo, error) {
	return n.api.Swarm().Peers(ctx)
}

// ID get self serviceNode info
func (n *nodeIPFS) ID(ctx context.Context) (pid *core.DataStoreNode, e error) {
	pid = &core.DataStoreNode{}
	e = n.api.Request("id").Exec(ctx, pid)
	if e != nil {
		return nil, e
	}
	return pid, nil
}

// PinAdd ...
func (n *nodeIPFS) PinAdd(ctx context.Context, hash string) (e error) {
	p := path.New(hash)
	return n.api.Pin().Add(ctx, p, options.Pin.Recursive(true))
}

// PinLS ...
func (n *nodeIPFS) PinLS(ctx context.Context) (pins []iface.Pin, e error) {
	return n.api.Pin().Ls(ctx, options.Pin.Type.Recursive())
}

// PinRm ...
func (n *nodeIPFS) PinRm(ctx context.Context, hash string) (e error) {
	p := path.New(hash)
	return n.api.Pin().Rm(ctx, p)
}

// IsReady ...
func (n *nodeIPFS) IsReady() bool {
	ma, err := multiaddr.NewMultiaddr(config.IPFSAddr())
	if err != nil {
		return false
	}
	api, e := httpapi.NewApi(ma)
	if e != nil {
		logE("new serviceNode ipfs", "error", e)
		return false
	}
	n.api = api
	return true
}

// Run ...
func (n *nodeIPFS) Run() {
	if n.lock.Load() {
		output("NodeClient is already running")
		return
	}
	n.lock.Store(true)
	defer n.lock.Store(false)
	if !n.IsReady() {
		output("waiting for ready")
		return
	}
	// get self serviceNode info
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Duration(n.cfg.IPFS.Timeout)*time.Second)
	defer cancelFunc()
	pid, e := n.ID(timeout)
	if e != nil {
		logE("run get id", "error", e)
		return
	}

	nid := pid.ID
	output("id", nid)
	// get ipfs swarm nodes
	publicNodes := make(map[string]bool)
	timeout2, cancelFunc2 := context.WithTimeout(context.Background(), time.Duration(n.cfg.IPFS.Timeout)*time.Second)
	defer cancelFunc2()
	infos, e := n.SwarmPeers(timeout2)
	if e != nil {
		logE("run get peers", "error", e)
		return
	}
	for _, info := range infos {
		output("peers", info.ID().String(), "ip", info.Address())
		conn, err := manet.Dial(info.Address())
		// p2p proxy serviceNode
		if err != nil {
			//TODO:
			output("err", err.Error())
			//ipfsAddr := "/ipfs/" + nodeID + "/p2p-circuit/ipfs/" + .ETHPeer
			//peers = append(peers, ipfsAddr)
		} else {
			ipfsAddr := info.Address().String()
			publicNodes[ipfsAddr] = true //append(publicNodes, ipfsAddr)
			conn.Close()
		}
	}
	//n.output("[当前IPFS总节点数]", len(peers)+len(publicNodes))
	output("exists IPFS nodes:", len(publicNodes))
	// get nodes info
	//if len(peers) == 0 {
	//	fmt.Println("<IPFS节点状态已是最新>")
	//	return
	//}
	cl := contract.Loader(n.cfg)
	err := cl.Node(func(node *node.AccelerateNode, opts *bind.TransactOpts) error {
		op := &bind.CallOpts{Pending: true}
		cPeers, err := node.GetIpfsNodes(op)
		if err != nil {
			log.Errorw("get ipfs node error", "tag", outputHead, "error", err)
			return err
		}
		cNodes, err := node.GetPublicIpfsNodes(op)
		if err != nil {
			log.Errorw("get public ipfs node error", "tag", outputHead, "error", err)
			return err
		}
		cPeers = decodeNodes(n.cfg, cPeers)
		cNodes = decodeNodes(n.cfg, cNodes)

		//TODO:fix sta
		//fmt.Println("[adding ipfs nodes]", difference(peers, cPeers))
		fmt.Println("[adding public ipfs nodes]", DiffStrArray(cNodes, publicNodes))
		// delete nodes
		var deleteIdx []int
		for _, dNode := range difference(cNodes, getAccessibleIPFSNodes(cNodes, "4001")) {
			for idx, cNode := range cNodes {
				if cNode == dNode {
					deleteIdx = append(deleteIdx, idx)
				}
			}
		}
		//TODO:fix end
		if len(deleteIdx) > 0 {
			var err error
			sort.Sort(sort.Reverse(sort.IntSlice(deleteIdx)))
			for _, idx := range deleteIdx {
				_, err = node.DeletePublicIpfsNodes(opts, uint32(idx))
			}

			if err != nil {
				fmt.Println("<删除失效节点失败>", err.Error())
			} else {
				fmt.Println("[删除失效节点成功]")
			}
		}

		// add new nodes
		//for _, n := range n.encodeNodes(difference(peers, cPeers)) {
		//	if n == "" {
		//		continue
		//	}
		//	_, err = ac.AddIpfsNodes(auth, []string{n})
		//}
		for _, n := range encodeNodes(n.cfg, DiffStrArray(cNodes, publicNodes)) {
			if n == "" {
				continue
			}
			_, err = node.AddPublicIpfsNodes(opts, []string{n})
		}

		if err != nil {
			fmt.Println("[添加节点失败]", err.Error())
		} else {
			fmt.Println("[添加节点成功] ")
		}
		return nil
	})

	if err != nil {
		log.Errorw("ipfs node process error", "tag", outputHead, "err", err)
		return
	}

	output("<IPFS同步完成>")
	return
}

// Address ...
type Address struct {
	Addr string
	Port string
}

// StringToAddr ...
func StringToAddr(s string) (a *Address, e error) {
	return
	//// consume trailing slashes
	//s = strings.TrimRight(s, "/")
	//
	//var b bytes.Buffer
	//sp := strings.Split(s, "/")
	//
	//if sp[0] != "" {
	//	return nil, fmt.Errorf("failed to parse multiaddr %q: must begin with /", s)
	//}
	//
	//// consume first empty elem
	//sp = sp[1:]
	//
	//if len(sp) == 0 {
	//	return nil, fmt.Errorf("failed to parse multiaddr %q: empty multiaddr", s)
	//}
	//
	//for len(sp) > 0 {
	//	log.Infow("sp info", "tag", outputHead, "info", sp)
	//	name := sp[0]
	//	p := multiaddr.ProtocolWithName(name)
	//	if p.Code == 0 {
	//		return nil, fmt.Errorf("failed to parse multiaddr %q: unknown protocol %s", s, sp[0])
	//	}
	//
	//	_, _ = b.Write(multiaddr.CodeToVarint(p.Code))
	//	log.Info(sp)
	//	sp = sp[1:]
	//	if p.Size == 0 { // no length.
	//		continue
	//	}
	//
	//	if len(sp) < 1 {
	//		return nil, fmt.Errorf("failed to parse multiaddr %q: unexpected end of multiaddr", s)
	//	}
	//
	//	if p.Path {
	//		// it's a path protocol (terminal).
	//		// consume the rest of the address as the next component.
	//		sp = []string{"/" + strings.Join(sp, "/")}
	//	}
	//
	//	a, err := p.Transcoder.StringToBytes(sp[0])
	//	if err != nil {
	//		return nil, fmt.Errorf("failed to parse multiaddr %q: invalid value %q for protocol %s: %s", s, sp[0], p.Name, err)
	//	}
	//	if p.Size < 0 { // varint size.
	//		_, _ = b.Write(multiaddr.CodeToVarint(len(a)))
	//	}
	//	b.Write(a)
	//	log.Info(sp)
	//	sp = sp[1:]
	//}
	//
	//return b.Bytes(), nil
}

// difference returns the elements in `a` that aren't in `b`.
func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func getAccessibleIPFSNodes(addresses []string, port string) []string {
	var accessible []string
	for _, address := range addresses {
		strs := strings.Split(address, "/ip4/")
		if len(strs) < 2 {
			continue
		}
		ip := strings.Split(strs[1], "/")[0]
		if len(ip) < 8 {
			continue
		}
		conn, err := net.DialTimeout("tcp", ip+":"+port, 3*time.Second)
		if err == nil {
			addr := strs[0] + "@" + ip + ":" + port
			accessible = append(accessible, addr)
			_ = conn.Close()
		} else {
			fmt.Println("[dial err]", err)
		}

	}
	return accessible
}
