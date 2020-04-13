package core

import (
	"go.uber.org/atomic"
	"sync"
)

// Node todo need fix
type _Node struct {
	RemoteAddr string
	Port       int
	Info       NodeInfo
	Hash       []string
}

// nodeStore ...
type nodeStore struct {
	nodes    sync.Map
	nodeSize *atomic.Int64
}

// NodeStore ...
type NodeStore interface {
	Add(info *NodeInfo)
	Check(key string) bool
	Get(key string) *NodeInfo
	Remove(key string)
	Length() int64
	Range(func(info *NodeInfo) bool)
}

// NewNodeStore ...
func NewNodeStore() NodeStore {
	return &nodeStore{
		nodes:    sync.Map{},
		nodeSize: atomic.NewInt64(0),
	}
}

// Remove ...
func (s *nodeStore) Remove(key string) {
	if s.Check(key) {
		s.nodeSize.Add(-1)
		s.nodes.Delete(key)
	}
}

// Add ...
func (s *nodeStore) Add(info *NodeInfo) {
	s.nodes.Store(info.Name, info)
	s.nodeSize.Add(1)
}

// Check ...
func (s *nodeStore) Check(key string) (b bool) {
	_, b = s.nodes.Load(key)
	return
}

// Get ...
func (s *nodeStore) Get(key string) *NodeInfo {
	if v, b := s.nodes.Load(key); b {
		return v.(*NodeInfo)
	}
	return nil
}

// Range ...
func (s *nodeStore) Range(f func(info *NodeInfo) bool) {
	s.nodes.Range(func(key, value interface{}) bool {
		return f(value.(*NodeInfo))
	})
}

// Length ...
func (s *nodeStore) Length() int64 {
	return s.nodeSize.Load()
}
