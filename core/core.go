package core

import "sync"

// NodeInfo ...
type NodeInfo struct {
	Name         string
	ContractAddr string
	DataAddr     string
	Version      string
}

// AccelerateCache ...
type AccelerateCache struct {
	nodes sync.Map
}

// Add ...
func (c *AccelerateCache) Add(info *NodeInfo) {
	c.nodes.Store(info.Name, info)
}

// Check ...
func (c *AccelerateCache) Check(key string) (b bool) {
	_, b = c.nodes.Load(key)
	return
}

// Get ...
func (c *AccelerateCache) Get(key string) *NodeInfo {
	if v, b := c.nodes.Load(key); b {
		return v.(*NodeInfo)
	}
	return nil
}

// Range ...
func (c *AccelerateCache) Range(f func(info *NodeInfo) bool) {
	c.nodes.Range(func(key, value interface{}) bool {
		return f(value.(*NodeInfo))
	})
}
