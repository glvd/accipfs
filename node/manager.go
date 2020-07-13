package node

import (
	"encoding/json"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/controller"
	"github.com/glvd/accipfs/core"
	"github.com/godcong/scdt"
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

// Manager ...
func Manager(cfg *config.Config, ctx *controller.APIContext) core.NodeManager {
	m := &manager{
		cfg:      cfg,
		initLoad: atomic.NewBool(false),
		path:     filepath.Join(cfg.Path, _nodes),
		expPath:  filepath.Join(cfg.Path, _expNodes),
		nodes:    nodeCacher(cfg),
		hashes:   hashCacher(cfg),
		t:        time.NewTicker(cfg.Node.BackupSeconds),
	}
	m.api = ctx.API(m)

	m.nodePool = mustPool(ants.DefaultAntsPoolSize, m.poolRun)

	//todo
	info, err := m.api.NodeAPI().NodeAddrInfo(&core.AddrReq{})
	if err != nil {
		return nil
	}
	m.local = core.NodeInfo{
		AddrInfo:        *info.AddrInfo,
		AgentVersion:    "",
		ProtocolVersion: "",
	}

	go m.loop()

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
	info, err := api.NodeAPI().NodeAddrInfo(&core.AddrReq{})
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
