package core

import (
	ma "github.com/multiformats/go-multiaddr"
)

// RecvCBFunc ...
type RecvCBFunc func(id string, v interface{}) ([]byte, error)

// Node ...
type Node interface {
	ID() string
	Addrs() []ma.Multiaddr
	DataStoreInfo() (DataStoreInfo, error)
	Info() (NodeInfo, error)
	Close() (err error)
	IsClosed() bool
	AppendAddr(addrs ...ma.Multiaddr)
	SendClose()
	Peers() ([]string, error)
	SendConnected() error
}
