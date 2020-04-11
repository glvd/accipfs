package service

import "github.com/glvd/accipfs/core"

type nodeManager struct {
	nodes map[string]core.NodeInfo
}
