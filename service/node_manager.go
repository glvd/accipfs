package service

import (
	"github.com/glvd/accipfs/config"
	"sync"
	"time"

	"github.com/glvd/accipfs/core"

	"github.com/gocacher/cacher"
	"go.uber.org/atomic"
)

// nodeManager ...
type nodeManager struct {
	cfg          *config.Config
	nodeCache    cacher.Cacher
	accountNodes sync.Map
	nodes        sync.Map
	faultNodes   sync.Map
	nodeSize     *atomic.Int64
}

// NodeManager ...
type NodeManager interface {
	Add(info *core.Node)
	Validate(key string, fs ...func(node *core.Node) bool) bool
	Get(key string) *core.Node
	Remove(key string)
	Length() int64
	Range(func(info *core.Node) bool)
	Fault(node *core.Node, fs ...func(info *core.Node) *core.Node)
	IsFault(key string) *core.Node
}

// NewNodeManager ...
func NewNodeManager(cfg *config.Config) NodeManager {
	return &nodeManager{
		cfg:      cfg,
		nodes:    sync.Map{},
		nodeSize: atomic.NewInt64(0),
	}
}

// Remove ...
func (s *nodeManager) Remove(key string) {
	s.nodeSize.Add(-1)
	s.nodes.Delete(key)
}

// Add ...
func (s *nodeManager) Add(info *core.Node) {
	if info.NodeType == core.NodeAccount {
		s.accountNodes.Store(info.Name, info)
	}
	s.nodes.Store(info.Name, info)
	s.nodeSize.Add(1)
}

// Validate ...
func (s *nodeManager) Validate(key string, fs ...func(node *core.Node) bool) (b bool) {
	n, exist := s.nodes.Load(key)
	if exist && fs != nil {
		node := n.(*core.Node)
		if b = fs[0](node); !b {
			s.Fault(node)
		}
	}
	return
}

// Fault ...
func (s *nodeManager) Fault(node *core.Node, fs ...func(info *core.Node) *core.Node) {
	s.Remove(node.NodeInfo.Name)
	if fs != nil {
		node = fs[0](node)
	}
	node.LastTime = time.Now()
	s.faultNodes.Store(node.NodeInfo.Name, node)
}

// IsFault ...
func (s *nodeManager) IsFault(key string) *core.Node {
	n, exist := s.nodes.Load(key)
	if exist {
		return n.(*core.Node)
	}
	return nil
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
