package core

import (
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

// RecvCBFunc ...
type RecvCBFunc func(id string, v interface{}) ([]byte, error)

// Node ...
type Node interface {
	ID() string
	Addrs() []ma.Multiaddr
	IPFSAddrInfo() (peer.AddrInfo, error)
	Info() (peer.AddrInfo, error)
	Close() (err error)
	IsClosed() bool
	AppendAddr(addrs ...ma.Multiaddr)
}
