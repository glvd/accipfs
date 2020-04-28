package cache

import (
	"encoding/json"
	"github.com/glvd/accipfs/config"
	"github.com/gocacher/cacher"
	"sync"
	"time"

	"github.com/glvd/accipfs/core"

	"go.uber.org/atomic"
)

// nodeCache ...
type nodeCache struct {
	cfg        *config.Config
	stop       *atomic.Bool
	pool       sync.Pool
	cache      cacher.Cacher
	nodes      sync.Map
	faultNodes sync.Map
	nodeSize   *atomic.Int64
}

// NodeCache ...
type NodeCache interface {
	Add(info *core.Node)
	Validate(key string, fs ...func(node *core.Node) bool)
	Get(key string) *core.Node
	Remove(key string)
	Length() int64
	Range(func(info *core.Node) bool)
	Fault(node *core.Node, fs ...func(info *core.Node))
	RecoveryFault(key string, fs ...func(info *core.Node)) (node *core.Node, ok bool)
}

// NewNodeCache ...
func NewNodeCache(cfg *config.Config) NodeCache {
	return &nodeCache{
		cfg:      cfg,
		stop:     atomic.NewBool(false),
		nodes:    sync.Map{},
		cache:    New(cfg),
		nodeSize: atomic.NewInt64(0),
	}
}

func (s *nodeCache) poolRun() {
	defer func() {
		if e := recover(); e != nil {
			logE("found error", "error", e)
		}
		if s.stop.Load() {

			return
		}
		go s.poolRun()
	}()
	for {
		if s.stop.Load() {
			return
		}
		if v := s.pool.Get(); v != nil {
			node := v.(*core.Node)
			if node.NodeType == core.NodeAccount {
				if err := cacher.Set(node.Name, marshalNode(node)); err != nil {
					panic(err)
				}
			}
			s.nodes.Store(node.Name, node)
			s.nodeSize.Add(1)
			continue
		}
		time.Sleep(3 * time.Second)
	}
}

// Remove ...
func (s *nodeCache) Remove(key string) {
	s.nodeSize.Add(-1)
	s.nodes.Delete(key)
}

// Add ...
func (s *nodeCache) Add(info *core.Node) {
	s.pool.Put(info)
}

// Validate ...
func (s *nodeCache) Validate(key string, fs ...func(node *core.Node) bool) {
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
func (s *nodeCache) Fault(node *core.Node, fs ...func(info *core.Node)) {
	s.Remove(node.NodeInfo.Name)
	if fs != nil {
		fs[0](node)
	}
	node.LastTime = time.Now()
	s.faultNodes.Store(node.NodeInfo.Name, node)
}

// RecoveryFault ...
func (s *nodeCache) RecoveryFault(key string, fs ...func(info *core.Node)) (node *core.Node, ok bool) {
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
func (s *nodeCache) LoadFault(key string) (*core.Node, bool) {
	n, exist := s.faultNodes.Load(key)
	if exist {
		return n.(*core.Node), exist
	}
	return nil, exist
}

// Get ...
func (s *nodeCache) Get(key string) *core.Node {
	if v, b := s.nodes.Load(key); b {
		return v.(*core.Node)
	}
	return nil
}

// GetAccount ...
func (s *nodeCache) GetAccount(key string) *core.Node {
	get, err := s.cache.Get(key)
	if err != nil {
		return nil
	}
	return unmarshalNode(get)
}

// Range ...
func (s *nodeCache) Range(f func(info *core.Node) bool) {
	s.nodes.Range(func(key, value interface{}) bool {
		return f(value.(*core.Node))
	})
}

// NodeHashes ...
func (s *nodeCache) NodeHashes(node *core.Node) []string {
	return nil
}

// HashNodes ...
func (s *nodeCache) HashNodes(hash string) []*core.Node {
	return nil
}

// Length ...
func (s *nodeCache) Length() int64 {
	return s.nodeSize.Load()
}

func marshalNode(node *core.Node) []byte {
	marshal, err := json.Marshal(node)
	if err != nil {
		panic(err)
	}
	return marshal
}

func unmarshalNode(b []byte) *core.Node {
	var node core.Node
	err := json.Unmarshal(b, &node)
	if err != nil {
		panic(err)
	}
	return &node
}
