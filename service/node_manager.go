package service

import "github.com/glvd/accipfs/core"

type nodeManager struct {
	nodes map[string]core.NodeInfo
}

// Add ...
func (m *nodeManager) Add(info *core.NodeInfo) {

}

// Get ...
func (m *nodeManager) Get(name string) *core.NodeInfo {
	return nil
}

// Filter ...
func (m *nodeManager) Filter(rule string) *core.NodeInfo {
	return nil
}
