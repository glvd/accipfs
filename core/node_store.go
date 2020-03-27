package core

import (
	"go.uber.org/atomic"
	"sync"
)

// Version ...
const Version = "0.0.1"

// DataStoreNode ...
type DataStoreNode struct {
	ID              string   `json:"ID"`
	PublicKey       string   `json:"PublicKey"`
	Addresses       []string `json:"Addresses"`
	AgentVersion    string   `json:"AgentVersion"`
	ProtocolVersion string   `json:"ProtocolVersion"`
}

// ContractNode ...
type ContractNode struct {
	Enode      string    `json:"enode"`
	Enr        string    `json:"enr"`
	ID         string    `json:"id"`
	IP         string    `json:"ip"`
	ListenAddr string    `json:"listenAddr"`
	Name       string    `json:"name"`
	Ports      Ports     `json:"ports"`
	Protocols  Protocols `json:"protocols"`
}

// Ports ...
type Ports struct {
	Discovery int64 `json:"discovery"`
	Listener  int64 `json:"listener"`
}

// Protocols ...
type Protocols struct {
	Eth Eth `json:"eth"`
}

// Eth ...
type Eth struct {
	Config     Config `json:"config"`
	Difficulty int64  `json:"difficulty"`
	Genesis    string `json:"genesis"`
	Head       string `json:"head"`
	Network    int64  `json:"network"`
}

// Config ...
type Config struct {
	ByzantiumBlock      int64  `json:"byzantiumBlock"`
	ChainID             int64  `json:"chainId"`
	Clique              Clique `json:"clique"`
	ConstantinopleBlock int64  `json:"constantinopleBlock"`
	Eip150Block         int64  `json:"eip150Block"`
	Eip150Hash          string `json:"eip150Hash"`
	Eip155Block         int64  `json:"eip155Block"`
	Eip158Block         int64  `json:"eip158Block"`
	HomesteadBlock      int64  `json:"homesteadBlock"`
}

// Clique ...
type Clique struct {
	Epoch  int64 `json:"epoch"`
	Period int64 `json:"period"`
}

// NodeInfo ...
type NodeInfo struct {
	Name       string
	RemoteAddr string
	Port       int
	Contract   ContractNode
	DataStore  DataStoreNode
	Version    string
}

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
