package service

import (
	"context"
	"fmt"
	"github.com/glvd/accipfs"
	"github.com/glvd/accipfs/config"
	"github.com/ipfs/go-ipfs-http-client"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"os/exec"
	"path/filepath"
)

const ipfsPath = ".ipfs"
const ipfsAPI = "api"

type nodeClientIPFS struct {
	cfg config.Config
	api *httpapi.HttpApi
}

type nodeServerIPFS struct {
	name string
	cmd  *exec.Cmd
}

// PeerID ...
type PeerID struct {
	Addresses       []string `json:"Addresses"`
	AgentVersion    string   `json:"AgentVersion"`
	ID              string   `json:"ID"`
	ProtocolVersion string   `json:"ProtocolVersion"`
	PublicKey       string   `json:"PublicKey"`
}

func newNodeIPFS(config config.Config) (*nodeClientIPFS, error) {
	api, e := httpapi.NewPathApi(filepath.Join(config.Path, ipfsPath))
	if e != nil {
		return nil, e
	}
	return &nodeClientIPFS{
		cfg: config,
		api: api,
	}, nil
}

// SwarmConnect ...
func (i *nodeClientIPFS) SwarmConnect(ctx context.Context, addr string) (e error) {
	ma, e := multiaddr.NewMultiaddr(addr)
	if e != nil {
		return e
	}
	info, e := peer.AddrInfoFromP2pAddr(ma)
	if e != nil {
		return e
	}
	e = i.api.Swarm().Connect(ctx, *info)
	if e != nil {
		return e
	}
	return nil
}

func (i *nodeClientIPFS) connect() (e error) {
	ma, err := multiaddr.NewMultiaddr(i.cfg.IPFS.Addr)
	if err != nil {
		return err
	}
	i.api, e = httpapi.NewApi(ma)
	return
}

// SwarmPeers ...
func (i *nodeClientIPFS) SwarmPeers(ctx context.Context) ([]iface.ConnectionInfo, error) {
	return i.api.Swarm().Peers(ctx)
}

// ID get self node info
func (i *nodeClientIPFS) ID(ctx context.Context) (pid *PeerID, e error) {
	pid = &PeerID{}
	e = i.api.Request("id").Exec(ctx, pid)
	if e != nil {
		return nil, e
	}
	return pid, nil
}

// PinAdd ...
func (i *nodeClientIPFS) PinAdd(ctx context.Context, hash string) (e error) {
	p := path.New(hash)
	return i.api.Pin().Add(ctx, p)
}

// PinLS ...
func (i *nodeClientIPFS) PinLS(ctx context.Context) (pins []iface.Pin, e error) {
	return i.api.Pin().Ls(ctx)
}

// PinRm ...
func (i *nodeClientIPFS) PinRm(ctx context.Context, hash string) (e error) {
	p := path.New(hash)
	return i.api.Pin().Rm(ctx, p)
}

// IsReady ...
func (i *nodeClientIPFS) IsReady() bool {
	//TODO
	return false
}

// Run ...
func (i *nodeClientIPFS) Run() {
	fmt.Println("<同步IPFS节点中>")
	if !i.IsReady() {
		fmt.Println("<waiting for ipfs ready>")
		return
	}
	//// get self node info
	//selfInfo, err := ipfs.ID()
	//if err != nil {
	//	fmt.Println("[获取本节点信息失败]", err.Error())
	//	return
	//}
	//nodeID := selfInfo.ID
	//// get ipfs swarm nodes
	//var peers []string
	//var publicNodes []string
	//resp := ipfs.SwarmPeers()
	//for _, peer := range resp.Peers {
	//	// check if ipfs node in public net
	//	ip, port := getAddressInfo(peer.Addr)
	//	conn, err := net.Dial("tcp", ip+":"+port)
	//	// p2p proxy node
	//	if err != nil {
	//		ipfsAddr := "/ipfs/" + nodeID + "/p2p-circuit/ipfs/" + peer.Peer
	//		peers = append(peers, ipfsAddr)
	//		// public node
	//	} else {
	//		ipfsAddr := peer.Addr + "/ipfs/" + peer.Peer
	//		publicNodes = append(publicNodes, ipfsAddr)
	//		conn.Close()
	//	}
	//}
	//fmt.Println("[当前IPFS总节点数]", len(peers)+len(publicNodes))
	//fmt.Println("[当前IPFS公网节点数]", len(publicNodes))
	//// get nodes info
	//if len(peers) == 0 {
	//	fmt.Println("<IPFS节点状态已是最新>")
	//	return
	//}
	//cl := eth.ContractLoader()
	//ac, auth, client := cl.AccelerateNode()
	//defer client.Close()
	//
	//cPeers, err := ac.GetIpfsNodes(nil)
	//cNodes, err := ac.GetPublicIpfsNodes(nil)
	//cPeers = decodeNodes(cPeers)
	//cNodes = decodeNodes(cNodes)
	//
	//fmt.Println("[adding ipfs nodes]", difference(peers, cPeers))
	//fmt.Println("[adding public ipfs nodes]", difference(publicNodes, cNodes))
	//// delete nodes
	//var deleteIdx []int
	//for _, dNode := range difference(cNodes, getAccessibleIpfsNodes(cNodes, "4001")) {
	//	for idx, cNode := range cNodes {
	//		if cNode == dNode {
	//			deleteIdx = append(deleteIdx, idx)
	//		}
	//	}
	//}
	//
	//if len(deleteIdx) > 0 {
	//	var err error
	//	sort.Sort(sort.Reverse(sort.IntSlice(deleteIdx)))
	//	for _, idx := range deleteIdx {
	//		_, err = ac.DeletePublicIpfsNodes(auth, uint32(idx))
	//	}
	//
	//	if err != nil {
	//		fmt.Println("<删除失效节点失败>", err.Error())
	//	} else {
	//		fmt.Println("[删除失效节点成功]")
	//	}
	//}
	//
	//// add new nodes
	//for _, n := range encodeNodes(difference(peers, cPeers)) {
	//	if n == "" {
	//		continue
	//	}
	//	_, err = ac.AddIpfsNodes(auth, []string{n})
	//}
	//for _, n := range encodeNodes(difference(publicNodes, cNodes)) {
	//	if n == "" {
	//		continue
	//	}
	//	_, err = ac.AddPublicIpfsNodes(auth, []string{n})
	//}
	//
	//if err != nil {
	//	fmt.Println("[添加节点失败]", err.Error())
	//} else {
	//	fmt.Println("[添加节点成功] ")
	//}
	fmt.Println("<IPFS同步完成>")
	return
}

// Start ...
func (n *nodeServerIPFS) Start() {
	fmt.Println("starting", n.name)
}

// NodeServerIPFS ...
func NodeServerIPFS(cfg config.Config) Node {
	cmd := exec.Command(cfg.IPFS.Name, "")
	cmd.Env = accipfs.Environ()
	return &nodeServerIPFS{cmd: cmd}
}
