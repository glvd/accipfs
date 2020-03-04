package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/glvd/accipfs"
	"github.com/glvd/accipfs/config"
	"github.com/goextension/log"
	"github.com/ipfs/go-ipfs-http-client"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"go.uber.org/atomic"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const ipfsPath = ".ipfs"
const ipfsAPI = "api"

type nodeClientIPFS struct {
	lock *atomic.Bool
	cfg  config.Config
	api  *httpapi.HttpApi
	out  *color.Color
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

	return &nodeClientIPFS{
		cfg: config,
		//api:  api,
		lock: atomic.NewBool(false),
		out:  color.New(color.FgBlue),
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
	api, e := httpapi.NewPathApi(filepath.Join(i.cfg.Path, ipfsPath))
	if e != nil {
		log.Errorf("new node: %w ", e)
		return false
	}
	i.api = api
	return true
}

func (i *nodeClientIPFS) output(v ...interface{}) {
	fmt.Print(outputHead, " ")
	fmt.Print("[IPFS]", " ")
	_, _ = i.out.Println(v...)
}

// Run ...
func (i *nodeClientIPFS) Run() {
	i.output("syncing node")
	if i.lock.Load() {
		i.output("ipfs node is already running")
		return
	}
	i.lock.Store(true)
	defer i.lock.Store(false)
	i.output("ipfs sync running")
	if !i.IsReady() {
		i.output("waiting for ready")
		return
	}
	//// get self node info
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Duration(i.cfg.IPFS.Timeout)*time.Second)
	defer cancelFunc()
	pid, e := i.ID(timeout)
	if e != nil {
		log.Errorw("run get id", "tag", outputHead, "error", e)
		return
	}

	nid := pid.ID
	i.output("id", nid)
	//// get ipfs swarm nodes
	//var peers []string
	//var publicNodes []string
	timeout2, cancelFunc2 := context.WithTimeout(context.Background(), time.Duration(i.cfg.IPFS.Timeout)*time.Second)
	defer cancelFunc2()
	infos, e := i.SwarmPeers(timeout2)
	if e != nil {
		log.Errorw(outputHead, "tag", "run get peers", "error", e)
		return
	}
	for _, info := range infos {
		i.output("peers", info.ID().String(), "ip", info.Address())
		info.Address()
		StringToAddr(info.Address().String())
	}
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
	i.output("<IPFS同步完成>")
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

// StringToAddr ...
func StringToAddr(s string) ([]byte, error) {
	// consume trailing slashes
	s = strings.TrimRight(s, "/")

	var b bytes.Buffer
	sp := strings.Split(s, "/")

	if sp[0] != "" {
		return nil, fmt.Errorf("failed to parse multiaddr %q: must begin with /", s)
	}

	// consume first empty elem
	sp = sp[1:]

	if len(sp) == 0 {
		return nil, fmt.Errorf("failed to parse multiaddr %q: empty multiaddr", s)
	}

	for len(sp) > 0 {
		log.Infow("sp info", "tag", outputHead, "info", sp)
		name := sp[0]
		p := multiaddr.ProtocolWithName(name)
		if p.Code == 0 {
			return nil, fmt.Errorf("failed to parse multiaddr %q: unknown protocol %s", s, sp[0])
		}
		_, _ = b.Write(multiaddr.CodeToVarint(p.Code))
		log.Info(sp)
		sp = sp[1:]
		if p.Size == 0 { // no length.
			continue
		}

		if len(sp) < 1 {
			return nil, fmt.Errorf("failed to parse multiaddr %q: unexpected end of multiaddr", s)
		}

		if p.Path {
			// it's a path protocol (terminal).
			// consume the rest of the address as the next component.
			sp = []string{"/" + strings.Join(sp, "/")}
		}

		a, err := p.Transcoder.StringToBytes(sp[0])
		if err != nil {
			return nil, fmt.Errorf("failed to parse multiaddr %q: invalid value %q for protocol %s: %s", s, sp[0], p.Name, err)
		}
		if p.Size < 0 { // varint size.
			_, _ = b.Write(multiaddr.CodeToVarint(len(a)))
		}
		b.Write(a)
		log.Info(sp)
		sp = sp[1:]
	}

	return b.Bytes(), nil
}
