package core

import (
	"net"
)

// LocalDataInfo ...
type LocalDataInfo struct {
	//ID   string
	Node NodeInfo
}

// NodeManager ...
type NodeManager interface {
	NodeAPI
	Local() LocalDataInfo
	Close()
	Push(n Node)
	Range(f func(key string, node Node) bool)
	Conn(c net.Conn) (Node, error)
	Store() error
	Load() error
}
