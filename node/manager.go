package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/godcong/scdt"
	ma "github.com/multiformats/go-multiaddr"
	mnet "github.com/multiformats/go-multiaddr-net"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/atomic"
	"net"
	"path/filepath"
	"sync"
	"time"
)

type manager struct {
	scdt.Listener
	initLoad        *atomic.Bool
	cfg             *config.Config
	t               *time.Ticker
	currentTS       int64
	ts              int64
	currentNodes    *atomic.Uint64
	path            string
	expPath         string
	api             core.API
	local           core.NodeInfo
	nodePool        *ants.PoolWithFunc
	connectNodes    sync.Map
	disconnectNodes sync.Map
	nodes           Cacher
	hashes          Cacher
}

var _nodes = "bl.nodes"
var _expNodes = "exp.nodes"
var _ core.NodeManager = &manager{}

// GlobalManager ...
var GlobalManager core.NodeManager

// InitManager ...
func InitManager(cfg *config.Config) core.NodeManager {
	if cfg.Node.BackupSeconds == 0 {
		cfg.Node.BackupSeconds = 30 * time.Second
	}

	m := &manager{
		cfg:      cfg,
		initLoad: atomic.NewBool(false),
		path:     filepath.Join(cfg.Path, _nodes),
		expPath:  filepath.Join(cfg.Path, _expNodes),
		nodes:    nodeCacher(cfg),
		hashes:   hashCacher(cfg),
		t:        time.NewTicker(cfg.Node.BackupSeconds),
	}
	m.nodePool = mustPool(ants.DefaultAntsPoolSize, m.poolRun)
	go m.loop()
	GlobalManager = m
	return m
}

// NodeAPI ...
func (m *manager) NodeAPI() core.NodeAPI {
	return m
}

func mustPool(size int, pf func(v interface{})) *ants.PoolWithFunc {
	withFunc, err := ants.NewPoolWithFunc(size, pf)
	if err != nil {
		panic(err)
	}
	return withFunc
}

// Store ...
func (m *manager) Store() (err error) {
	m.connectNodes.Range(func(key, value interface{}) bool {
		keyk, keyb := key.(string)
		valv, valb := value.(core.Marshaler)
		if !valb || !keyb {
			return true
		}
		err := m.nodes.Store(keyk, valv)
		if err != nil {
			return false
		}
		return true
	})
	//return the last err
	return
}

// Link ...
func (m *manager) Link(req *core.NodeLinkReq) (*core.NodeLinkResp, error) {
	fmt.Printf("connect info:%+v\n", req.Addrs)
	for _, addr := range req.Addrs {
		multiaddr, err := ma.NewMultiaddr(addr)
		if err != nil {
			fmt.Printf("parse addr(%v) failed(%v)\n", addr, err)
			continue
		}
		dial, err := mnet.Dial(multiaddr)
		if err != nil {
			fmt.Printf("link failed(%v)\n", err)
			continue
		}
		conn, err := m.Conn(dial)
		if err != nil {
			return nil, err
		}
		info, err := conn.Info()
		if err != nil {
			return nil, err
		}
		return &core.NodeLinkResp{
			NodeInfo: info,
		}, nil
	}
	return &core.NodeLinkResp{}, errors.New("all request was failed")
}

// Unlink ...
func (m *manager) Unlink(req *core.NodeUnlinkReq) (*core.NodeUnlinkResp, error) {
	panic("implement me")
}

// NodeAddrInfo ...
func (m *manager) NodeAddrInfo(req *core.AddrReq) (*core.AddrResp, error) {
	load, ok := m.connectNodes.Load(req.ID)
	if !ok {
		return &core.AddrResp{}, fmt.Errorf("node not found id(%s)", req.ID)
	}
	v, b := load.(core.Node)
	if !b {
		return &core.AddrResp{}, fmt.Errorf("transfer to node failed id(%s)", req.ID)
	}
	panic("//todo")
	fmt.Print(v)
	return nil, nil
}

// List ...
func (m *manager) List(req *core.NodeListReq) (*core.NodeListResp, error) {
	nodes := make(map[string]core.NodeInfo)
	m.Range(func(key string, node core.Node) bool {
		info, err := node.Info()
		if err != nil {
			return true
		}
		nodes[key] = info
		return true
	})
	return &core.NodeListResp{Nodes: nodes}, nil
}

// Load ...
func (m *manager) Load() error {
	m.nodes.Range(func(hash string, value string) bool {
		var addrInfo core.AddrInfo
		err := json.Unmarshal([]byte(value), &addrInfo)
		if err != nil {
			log.Errorw("load addr info failed", "err", err)
		}
		for multiaddr := range addrInfo.Addrs {
			connectNode, err := ConnectNode(multiaddr, 0, m.local, m.api)
			if err != nil {
				continue
			}
			m.connectNodes.Store(hash, connectNode)
			return true
		}
		return true
	})
	return nil
}

// StateEx State Examination checks the node status
func (m *manager) StateEx(id string, f func(node core.Node) bool) {
	if f == nil {
		return
	}
	node, ok := m.connectNodes.Load(id)
	if ok {
		if f(node.(core.Node)) {
			m.connectNodes.Delete(id)
			m.disconnectNodes.Store(id, node)
		}
	}

	exp, ok := m.disconnectNodes.Load(id)
	if ok {
		if f(exp.(core.Node)) {
			m.disconnectNodes.Delete(id)
			m.connectNodes.Store(id, exp)
		}
	}
}

// Range ...
func (m *manager) Range(f func(key string, node core.Node) bool) {
	m.connectNodes.Range(func(key, value interface{}) bool {
		k, b1 := key.(string)
		n, b2 := value.(core.Node)
		if !b1 || !b2 {
			return true
		}
		if f != nil {
			return f(k, n)
		}
		return false
	})
}

// Push ...
func (m *manager) Push(node core.Node) {
	m.ts = time.Now().Unix()
	m.connectNodes.Store(node.ID(), node)
}

// save nodes
func (m *manager) loop() {
	if m.initLoad.Load() {
		err := m.Load()
		if err != nil {
			log.Errorw("load node failed", "err", err)
		}
	}
	for {
		<-m.t.C
		if m.ts != m.currentTS {
			if err := m.Store(); err != nil {
				continue
			}
			m.currentTS = m.ts
		}
	}
}

// Conn ...
func (m *manager) Conn(c net.Conn) (core.Node, error) {
	acceptNode, err := CoreNode(c, m.local, m.api)
	if err != nil {
		return nil, err
	}

	m.nodePool.Invoke(acceptNode)
	return acceptNode, nil
}

func (m *manager) poolRun(v interface{}) {
	n, b := v.(core.Node)
	if !b {
		return
	}
	store, loaded := m.connectNodes.Load(n.ID())
	if loaded {
		nbase := store.(core.Node)
		if n.Addrs() != nil {
			nbase.AppendAddr(n.Addrs()...)
		}
		n.Close()
		return
	}
	if !n.IsClosed() {
		m.Push(n)
	}

}

func decodeNode(m core.NodeManager, b []byte, api core.API) error {
	nodes := map[string]jsonNode{}
	err := json.Unmarshal(b, &nodes)
	if err != nil {
		return err
	}
	info, err := m.NodeAPI().NodeAddrInfo(&core.AddrReq{})
	if err != nil {
		return err
	}
	for _, nodes := range nodes {
		for _, addr := range nodes.Addrs {
			connectNode, err := ConnectNode(addr, 0, core.NodeInfo{
				AddrInfo:        *info.AddrInfo,
				AgentVersion:    "",
				ProtocolVersion: "",
			}, api)
			if err != nil {
				continue
			}
			m.Push(connectNode)
			break
		}
	}
	return nil
}

func encodeNode(node core.Node) ([]byte, error) {
	n := map[string]jsonNode{
		node.ID(): {Addrs: node.Addrs()},
	}
	return json.Marshal(n)
}

// Close ...
func (m *manager) Close() {
	m.hashes.Close()
}
