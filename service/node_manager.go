package service

import (
	"github.com/glvd/accipfs/core"
	"go.uber.org/atomic"
	"sync"
)

// nodeManager ...
type nodeManager struct {
	nodes    sync.Map
	dummy    sync.Map
	nodeSize *atomic.Int64
}

// NodeManager ...
type NodeManager interface {
	Add(info *core.Node)
	Check(key string) bool
	Get(key string) *core.Node
	Remove(key string)
	Length() int64
	Range(func(info *core.Node) bool)
}

// NewNodeManager ...
func NewNodeManager() *nodeManager {
	return &nodeManager{
		nodes:    sync.Map{},
		nodeSize: atomic.NewInt64(0),
	}
}

// Remove ...
func (s *nodeManager) Remove(key string) {
	if s.Check(key) {
		s.nodeSize.Add(-1)
		s.nodes.Delete(key)
	}
}

// Add ...
func (s *nodeManager) Add(info *core.Node) {
	s.nodes.Store(info.Name, info)
	s.nodeSize.Add(1)
}

// Check ...
func (s *nodeManager) Check(key string) (b bool) {
	_, b = s.nodes.Load(key)
	return
}

// Get ...
func (s *nodeManager) Get(key string) *core.Node {
	if v, b := s.nodes.Load(key); b {
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
