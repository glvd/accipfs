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
	path            string
	expPath         string
	local           core.SafeLocalData
	nodePool        *ants.PoolWithFunc
	currentNodes    *atomic.Int32
	connectNodes    sync.Map
	disconnectNodes sync.Map
	nodes           Cacher
	hashes          Cacher
	RequestLD       func() ([]string, error)
	gc              *atomic.Bool
}

var _nodes = "bl.nodes"
var _expNodes = "exp.nodes"
var _ core.NodeManager = &manager{}

// InitManager ...
func InitManager(cfg *config.Config) (core.NodeManager, error) {
	if cfg.Node.BackupSeconds == 0 {
		cfg.Node.BackupSeconds = 30 * time.Second
	}
	data := core.DefaultLocalData()
	m := &manager{
		cfg:      cfg,
		initLoad: atomic.NewBool(false),
		path:     filepath.Join(cfg.Path, _nodes),
		expPath:  filepath.Join(cfg.Path, _expNodes),
		nodes:    nodeCacher(cfg),
		hashes:   hashCacher(cfg),
		local:    data.Safe(),
		t:        time.NewTicker(cfg.Node.BackupSeconds),
	}
	m.nodePool = mustPool(cfg.Node.PoolMax, m.poolRun)
	go m.loop()

	return m, nil
}

// NodeAPI ...
func (m *manager) NodeAPI() core.NodeAPI {
	return m
}

func mustPool(size int, pf func(v interface{})) *ants.PoolWithFunc {
	if size == 0 {
		size = ants.DefaultAntsPoolSize
	}
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
	if !m.local.Data().Initialized {
		return &core.NodeLinkResp{}, errors.New("you are not ready for connection")
	}
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
		conn, err := m.newConn(dial)
		if err != nil {
			return &core.NodeLinkResp{}, err
		}
		id := conn.ID()
		getNode, b := m.GetNode(id)
		if b {
			conn = getNode
		} else {
			//use conn
		}
		info, err := conn.GetInfo()
		if err != nil {
			return &core.NodeLinkResp{}, err
		}
		return &core.NodeLinkResp{
			NodeInfo: info,
		}, nil
	}
	return &core.NodeLinkResp{}, errors.New("all request was failed")
}

// Unlink ...
func (m *manager) Unlink(req *core.NodeUnlinkReq) (*core.NodeUnlinkResp, error) {
	if len(req.Peers) == 0 {
		return &core.NodeUnlinkResp{}, nil
	}
	for i := range req.Peers {
		m.connectNodes.Delete(req.Peers[i])
	}
	return &core.NodeUnlinkResp{}, nil
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
		info, err := node.GetInfo()
		if err != nil {
			return true
		}
		nodes[key] = info
		return true
	})
	return &core.NodeListResp{Nodes: nodes}, nil
}

// Local ...
func (m *manager) Local() core.SafeLocalData {
	return m.local
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
			connectNode, err := ConnectNode(multiaddr, 0, m.local)
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

// newConn ...
func (m *manager) newConn(c net.Conn) (core.Node, error) {
	acceptNode, err := CoreNode(c, m.local)
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
	pushed := false
	defer func() {
		fmt.Println("id", n.ID(), "was exit")
		//wait client close itself
		time.Sleep(500 * time.Millisecond)
		n.Close()
		if pushed {
			m.connectNodes.Delete(n.ID())
		}
	}()
	id := n.ID()
	fmt.Println("user connect:", id)
	if id == "" {
		//wait client get base info
		time.Sleep(3 * time.Second)
		_ = n.SendConnected()
		return
	}
	old, loaded := m.connectNodes.Load(id)
	log.Infow("new connection", "new", n.ID(), "isload", loaded)
	if loaded {
		nbase := old.(core.Node)
		if n.Addrs() != nil {
			nbase.AppendAddr(n.Addrs()...)
		}
		//wait client get base info
		time.Sleep(3 * time.Second)
		_ = n.SendConnected()
		return
	}

	if !n.IsClosed() {
		fmt.Println("node added:", n.ID())
		pushed = true
		m.Push(n)
	}

	info, err := n.GetInfo()
	if err == nil {
		m.local.Update(func(data *core.LocalData) {
			data.Nodes[info.ID] = info
		})
	}
	for !n.IsClosed() {
		peers, err := n.Peers()
		for _, peer := range peers {
			if len(peer.Addrs) != 0 {
				m.connectMultiAddrs(peer)
			}
		}
		lds, err := n.LDs()
		if err != nil {
			fmt.Println("failed to get link data", err)
			if err == ErrNoData {
				continue
			}
			return
		}
		for _, ld := range lds {
			fmt.Println("from:", n.ID(), "list:", ld)
		}
		time.Sleep(5 * time.Second)
	}
}

func decodeNode(m core.NodeManager, b []byte, api core.API) error {
	nodes := map[string]jsonNode{}
	err := json.Unmarshal(b, &nodes)
	if err != nil {
		return err
	}

	for _, nodes := range nodes {
		for _, addr := range nodes.Addrs {
			connectNode, err := ConnectNode(addr, 0, m.Local())
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

// GetNode ...
func (m *manager) GetNode(id string) (n core.Node, b bool) {
	load, ok := m.connectNodes.Load(id)
	if ok {
		n, b = load.(core.Node)
		return
	}
	return
}

// AllNodes ...
func (m *manager) AllNodes() (map[string]core.Node, int, error) {
	nodes := make(map[string]core.Node)
	count := 0
	m.Range(func(key string, node core.Node) bool {
		nodes[key] = node
		node.Addrs()
		count++
		return true
	})
	return nodes, count, nil
}

// Add ...
func (m *manager) Add(req *core.AddReq) (*core.AddResp, error) {
	m.local.Update(func(data *core.LocalData) {
		data.LDs[req.Hash] = 0
	})
	return &core.AddResp{
		IsSuccess: true,
		Hash:      req.Hash,
	}, nil
}

// Conn ...
func (m *manager) Conn(c net.Conn) (core.Node, error) {
	return m.newConn(c)
}

func (m *manager) addNode(n core.Node) (bool, error) {
	if m.currentNodes.Load() > int32(m.cfg.Node.ConnectMax) {
		go m.nodeGC()
	}
	m.Push(n)
	return true, nil
}

func (m *manager) nodeGC() {
	if !m.gc.CAS(false, true) {
		return
	}
	defer m.gc.Store(false)
	m.connectNodes.Range(func(key, value interface{}) bool {
		v, b := value.(core.Node)
		if !b {
			m.connectNodes.Delete(key)
			return true
		}
		ping, err := v.Ping()
		if err != nil {
			m.connectNodes.Delete(key)
			return true
		}
		if ping != "pong" {
			m.connectNodes.Delete(key)
			return true
		}
		return true
	})

}

func (m *manager) connectMultiAddrs(info core.NodeInfo) {
	if info.ID == m.cfg.Identity {
		return
	}
	_, ok := m.connectNodes.Load(info.ID)
	if ok {
		return
	}
	addrs := info.GetAddrs()
	if addrs == nil {
		return
	}
	for _, addr := range info.GetAddrs() {
		dial, err := mnet.Dial(addr)
		if err != nil {
			fmt.Printf("link failed(%v)\n", err)
			continue
		}
		conn, err := m.newConn(dial)
		if err != nil {
			return
		}
		id := conn.ID()
		getNode, b := m.GetNode(id)
		if b {
			conn = getNode
		} else {
			//use conn
		}
	}
}
