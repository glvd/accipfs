package cache

import (
	"github.com/glvd/accipfs/config"
	"github.com/gocacher/cacher"
	"sync"
	"time"

	"github.com/glvd/accipfs/core"

	"go.uber.org/atomic"
)

// nodeManager ...
type nodeManager struct {
	cfg          *config.Config
	stop         *atomic.Bool
	pool         sync.Pool
	cache        cacher.Cacher
	accountNodes sync.Map
	nodes        sync.Map
	faultNodes   sync.Map
	nodeSize     *atomic.Int64
}

// NodeManager ...
type NodeManager interface {
	Add(info *core.Node)
	Validate(key string, fs ...func(node *core.Node) bool)
	Get(key string) *core.Node
	Remove(key string)
	Length() int64
	Range(func(info *core.Node) bool)
	Fault(node *core.Node, fs ...func(info *core.Node))
	RecoveryFault(key string, fs ...func(info *core.Node)) (node *core.Node, ok bool)
}

// NewNodeManager ...
func NewNodeManager(cfg *config.Config) NodeManager {
	return &nodeManager{
		cfg:      cfg,
		stop:     atomic.NewBool(false),
		nodes:    sync.Map{},
		nodeSize: atomic.NewInt64(0),
	}
}

func (s *nodeManager) poolRun() {
	for {
		if s.stop.Load() {
			return
		}
		if v := s.pool.Get(); v != nil {
			node := v.(*core.Node)
			if node.NodeType == core.NodeAccount {
				s.accountNodes.Store(node.Name, node)
			}
			s.nodeSize.Add(1)
			s.nodes.Store(node.Name, node)
			continue
		}
		time.Sleep(3 * time.Second)
	}
}

// Remove ...
func (s *nodeManager) Remove(key string) {
	s.nodeSize.Add(-1)
	s.nodes.Delete(key)
}

// Add ...
func (s *nodeManager) Add(info *core.Node) {
	s.pool.Put(info)
}

// Validate ...
func (s *nodeManager) Validate(key string, fs ...func(node *core.Node) bool) {
	n, exist := s.nodes.Load(key)
	if exist && fs != nil {
		node := n.(*core.Node)
		if b := fs[0](node); !b {
			s.Fault(node)
		}
	}
	return
}

// Fault ...
func (s *nodeManager) Fault(node *core.Node, fs ...func(info *core.Node)) {
	s.Remove(node.NodeInfo.Name)
	if fs != nil {
		fs[0](node)
	}
	node.LastTime = time.Now()
	s.faultNodes.Store(node.NodeInfo.Name, node)
}

// RecoveryFault ...
func (s *nodeManager) RecoveryFault(key string, fs ...func(info *core.Node)) (node *core.Node, ok bool) {
	load, ok := s.faultNodes.Load(key)
	if !ok {
		return
	}
	node = load.(*core.Node)
	for _, f := range fs {
		f(node)
	}
	s.Add(node)
	return node, true
}

// IsFault ...
func (s *nodeManager) LoadFault(key string) (*core.Node, bool) {
	n, exist := s.faultNodes.Load(key)
	if exist {
		return n.(*core.Node), exist
	}
	return nil, exist
}

// Get ...
func (s *nodeManager) Get(key string) *core.Node {
	if v, b := s.nodes.Load(key); b {
		return v.(*core.Node)
	}
	return nil
}

// GetAccount ...
func (s *nodeManager) GetAccount(key string) *core.Node {
	if v, b := s.accountNodes.Load(key); b {
		return v.(*core.Node)
	}
	return nil
}

// Range ...
func (s *nodeManager) Range(f func(info *core.Node) bool) {
	s.nodes.Range(func(key, value interface{}) bool {
		return f(value.(*core.Node))
	})
}

// Length ...
func (s *nodeManager) Length() int64 {
	return s.nodeSize.Load()
}

func faultTimeCheck(fault *core.Node, limit int64) (remain int64, fa bool) {
	now := time.Now().Unix()
	f := fault.LastTime.Unix() + limit
	remain = -(now - f)
	if remain < 0 {
		remain = 0
	}
	return remain, remain <= 0
}
