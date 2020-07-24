package node

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"sync"
	"time"

	"github.com/glvd/accipfs/basis"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/godcong/scdt"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	mnet "github.com/multiformats/go-multiaddr-net"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/atomic"
)

type manager struct {
	scdt.Listener
	//initLoad        *atomic.Bool
	loopOnce        *sync.Once
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
	nodes           Cacher //all node caches
	hashNodes       Cacher //hash cache nodes
	RequestLD       func() ([]string, error)
	gc              *atomic.Bool
	addrCB          func(info peer.AddrInfo) error
}

var _nodes = "bl.nodes"
var _expNodes = "exp.nodes"
var _ core.NodeManager = &manager{}

// InitManager ...
func InitManager(cfg *config.Config) (core.NodeManager, error) {
	if cfg.Node.BackupSeconds == 0 {
		cfg.Node.BackupSeconds = 30
	}
	data := core.DefaultLocalData()
	m := &manager{
		cfg:      cfg,
		loopOnce: &sync.Once{},
		//initLoad:  atomic.NewBool(false),
		path:      filepath.Join(cfg.Path, _nodes),
		expPath:   filepath.Join(cfg.Path, _expNodes),
		nodes:     nodeCacher(cfg),
		hashNodes: hashCacher(cfg),
		local:     data.Safe(),
		t:         time.NewTicker(cfg.Node.BackupSeconds * time.Second),
	}
	m.nodePool = mustPool(cfg.Node.PoolMax, m.mainProc)
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
		log.Infow("store", "key", keyk, "keyOK", keyb, "value", valv, "valOK", valb)
		if !valb || !keyb {
			return true
		}
		err := m.nodes.Store(keyk, valv)
		if err != nil {
			return false
		}
		fmt.Println("node", keyk, "was stored")
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
		if req.Timeout == 0 {
			req.Timeout = 5 * time.Second
		}
		d := mnet.Dialer{
			Dialer: net.Dialer{
				Timeout: req.Timeout,
			},
		}
		dial, err := d.Dial(multiaddr)
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
	//todo:need optimization
	//nodes := make(map[string]core.NodeInfo)
	//m.Range(func(key string, node core.Node) bool {
	//	info, err := node.GetInfo()
	//	if err != nil {
	//		return true
	//	}
	//	nodes[key] = info
	//	return true
	//})
	return &core.NodeListResp{Nodes: m.local.Data().Nodes}, nil
}

// Local ...
func (m *manager) Local() core.SafeLocalData {
	return m.local
}

// Load ...
func (m *manager) Load() error {
	m.nodes.Range(func(hash string, value string) bool {
		log.Infow("range node", "hash", hash, "value", value)
		var ninfo core.NodeInfo
		err := ninfo.Unmarshal([]byte(value))
		//err := json.Unmarshal([]byte(value), &ninfo)
		if err != nil {
			log.Errorw("load addr info failed", "err", err)
			return true
		}
		for multiaddr := range ninfo.Addrs {
			fmt.Println("load node from address:", multiaddr.String())
			connectNode, err := ConnectNode(multiaddr, 0, m.local)
			if err != nil {
				continue
			}
			m.connectNodes.Store(hash, connectNode)
			return true
		}
		return true
	})
	//start loop after first load
	m.loopOnce.Do(func() {
		go m.loop()
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
	//if m.initLoad.CAS(false, true) {
	//	err := m.Load()
	//	if err != nil {
	//		log.Errorw("load node failed", "err", err)
	//	}
	//}
	for {
		<-m.t.C
		fmt.Println("store new node")
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

func (m *manager) mainProc(v interface{}) {
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
		//wait remote client get base info
		time.Sleep(3 * time.Second)
		_ = n.SendConnected()
		return
	}
	old, loaded := m.connectNodes.Load(id)
	log.Infow("new connection", "new", n.ID(), "isload", loaded)
	if loaded {
		nbase := old.(core.Node)
		if !nbase.IsClosed() {
			if n.Addrs() != nil {
				nbase.AppendAddr(n.Addrs()...)
			}
			//wait remote client get base info
			time.Sleep(3 * time.Second)
			_ = n.SendConnected()
			return
		}
	}
	//get remote node info
	info, err := n.GetInfo()
	if err == nil {
		if m.local.Data().Node.ID != info.ID {
			m.local.Update(func(data *core.LocalData) {
				data.Nodes[info.ID] = info
			})
			//err := m.nodes.Store(info.ID, info)
			//if err != nil {
			//	log.Errorw("sotre nodes failed", "err", err)
			//}
			m.connectRemoteDataStore(info.DataStore)
		}
	}

	if !n.IsClosed() {
		fmt.Println("node added:", n.ID())
		pushed = true
		m.Push(n)
	}

	for !n.IsClosed() {
		peerDone := m.syncPeers(n)
		lds, err := n.LDs()
		if err != nil {
			fmt.Println("failed to get link data", err)
			if err == ErrNoData {
				continue
			}
			//close connect when found err?
			return
		}
		for _, ld := range lds {
			err := m.hashNodes.Update(ld, func(bytes []byte) (core.Marshaler, error) {
				nodes := NewNodes()
				err := nodes.Unmarshal(bytes)
				if err != nil {
					return nil, err
				}
				nodes.n[n.ID()] = true
				return nodes, nil
			})
			if err != nil {
				continue
			}
			fmt.Println("from:", n.ID(), "list:", ld)
		}
		//wait something done
		<-peerDone
		time.Sleep(5 * time.Second)
	}
}

func (m *manager) connectRemoteDataStore(info core.DataStoreInfo) {
	timeout, cancelFunc := context.WithTimeout(context.TODO(), time.Second*30)
	defer cancelFunc()

	if m.addrCB != nil {
		addresses, err := basis.ParseAddresses(timeout, info.Addresses)
		if err != nil {
			return
		}
		total := len(addresses)
		for _, addr := range addresses {
			err = m.addrCB(addr)
			if err != nil {
				total = total - 1
				continue
			}
		}
		if total <= 0 {
			log.Infow("addr callback failed", "err", err, "addrinfo", info.Addresses)
		}
	}

}

func (m *manager) syncPeers(n core.Node) <-chan bool {
	done := make(chan bool)
	go func() {
		peers, err := n.Peers()
		if err != nil {
			done <- false
		}
		for _, peer := range peers {
			if len(peer.Addrs) != 0 {
				m.connectMultiAddr(peer)
			}
		}
		done <- true
	}()
	return done
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
	m.nodes.Close()
	m.hashNodes.Close()
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
	m.connectNodes.Range(func(key, value interface{}) (deleted bool) {
		defer func() {
			keyStr, b := key.(string)
			if !b {
				return
			}
			if deleted {
				m.local.Update(func(data *core.LocalData) {
					delete(data.Nodes, keyStr)
				})
			}
		}()
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

func (m *manager) connectMultiAddr(info core.NodeInfo) error {
	if info.ID == m.cfg.Identity {
		return nil
	}
	_, ok := m.connectNodes.Load(info.ID)
	if ok {
		return nil
	}
	addrs := info.GetAddrs()
	if addrs == nil {
		return nil
	}
	for _, addr := range addrs {
		dialer := mnet.Dialer{
			Dialer: net.Dialer{
				Timeout: 3 * time.Second,
			},
		}
		dial, err := dialer.Dial(addr)
		if err != nil {
			fmt.Printf("link failed(%v)\n", err)
			return err
		}
		fmt.Printf("link success(%v)\n", addr)
		_, err = m.newConn(dial)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("no link connect")
}

// RegisterAddrCallback ...
func (m *manager) RegisterAddrCallback(f func(info peer.AddrInfo) error) {
	m.addrCB = f
}

// ConnRemoteFromHash ...
func (m *manager) ConnRemoteFromHash(hash string) error {
	var nodes Nodes
	err := m.hashNodes.Load(hash, &nodes)
	if err != nil {
		return err
	}
	for s := range nodes.n {
		addr, err := DialFromStringAddr(s, 0)
		if err != nil {
			continue
		}
		_, err = m.Conn(addr)
		if err != nil {
			continue
		}
	}
	return nil
}
