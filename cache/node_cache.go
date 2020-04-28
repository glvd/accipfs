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
	n := &nodeCache{
		cfg:      cfg,
		stop:     atomic.NewBool(false),
		nodes:    sync.Map{},
		cache:    New(cfg),
		nodeSize: atomic.NewInt64(0),
	}
	go n.poolRun()
	return n
}

func (c *nodeCache) poolRun() {
	defer func() {
		if e := recover(); e != nil {
			logE("found error", "error", e)
		}
		if c.stop.Load() {
			return
		}
		go c.poolRun()
	}()
	for {
		if c.stop.Load() {
			return
		}
		if v := c.pool.Get(); v != nil {
			node := v.(*core.Node)
			if node.NodeType == core.NodeAccount {
				if err := cacher.Set(node.Name, marshalNode(node)); err != nil {
					panic(err)
				}
			}
			c.nodes.Store(node.Name, node)
			c.nodeSize.Add(1)
			continue
		}
		time.Sleep(3 * time.Second)
	}
}

// Remove ...
func (c *nodeCache) Remove(key string) {
	c.nodeSize.Add(-1)
	c.nodes.Delete(key)
}

// Add ...
func (c *nodeCache) Add(info *core.Node) {
	c.pool.Put(info)
}

// Validate ...
func (c *nodeCache) Validate(key string, fs ...func(node *core.Node) bool) {
	n, exist := c.nodes.Load(key)
	if exist && fs != nil {
		node := n.(*core.Node)
		if b := fs[0](node); !b {
			c.Fault(node)
		}
	}
	return
}

// Fault ...
func (c *nodeCache) Fault(node *core.Node, fs ...func(info *core.Node)) {
	c.Remove(node.NodeInfo.Name)
	if fs != nil {
		fs[0](node)
	}
	node.LastTime = time.Now()
	c.faultNodes.Store(node.NodeInfo.Name, node)
}

// RecoveryFault ...
func (c *nodeCache) RecoveryFault(key string, fs ...func(info *core.Node)) (node *core.Node, ok bool) {
	load, ok := c.faultNodes.Load(key)
	if !ok {
		return
	}
	node = load.(*core.Node)
	for _, f := range fs {
		f(node)
	}
	c.Add(node)
	return node, true
}

// IsFault ...
func (c *nodeCache) LoadFault(key string) (*core.Node, bool) {
	n, exist := c.faultNodes.Load(key)
	if exist {
		return n.(*core.Node), exist
	}
	return nil, exist
}

// Get ...
func (c *nodeCache) Get(key string) *core.Node {
	if v, b := c.nodes.Load(key); b {
		return v.(*core.Node)
	}
	return nil
}

// GetAccount ...
func (c *nodeCache) GetAccount(key string) *core.Node {
	get, err := c.cache.Get(key)
	if err != nil {
		return nil
	}
	return unmarshalNode(get)
}

// Range ...
func (c *nodeCache) Range(f func(info *core.Node) bool) {
	c.nodes.Range(func(key, value interface{}) bool {
		return f(value.(*core.Node))
	})
}

// NodeHashes ...
func (c *nodeCache) NodeHashes(node *core.Node) []string {
	return nil
}

// Length ...
func (c *nodeCache) Length() int64 {
	return c.nodeSize.Load()
}

// Cancel ...
func (c *nodeCache) Cancel() {
	c.stop.Store(true)
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
