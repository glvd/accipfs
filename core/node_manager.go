package core

import (
	"net"
)

// NodeManager ...
type NodeManager interface {
	NodeAPI() NodeAPI
	Local() SafeLocalData
	Close()
	Push(n Node)
	Range(f func(key string, node Node) bool)
	Conn(c net.Conn) (Node, error)
	Store() error
	Load() error
	RegisterLDRequest(func() ([]string, error))
	RegisterAddrRequest(f func() ([]string, error))
}
