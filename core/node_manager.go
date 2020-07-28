package core

import (
	"github.com/libp2p/go-libp2p-core/peer"
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
	Save() error
	Load() error

	//RegisterLDRequest(func() ([]string, error))
	RegisterAddrCallback(f func(info peer.AddrInfo) error)
	ConnRemoteFromHash(hash string) error
}
