package core

import (
	ma "github.com/multiformats/go-multiaddr"
)

// RecvCBFunc ...
type RecvCBFunc func(id string, v interface{}) ([]byte, error)

// Node ...
type Node interface {
	ID() string
	Ping() (string, error)
	Addrs() []ma.Multiaddr
	DataStoreInfo() (DataStoreInfo, error)
	GetInfo() (NodeInfo, error)
	Close() (err error)
	IsClosed() bool
	AppendAddr(addrs ...ma.Multiaddr)
	SendClose()
	Peers() ([]NodeInfo, error)
	SendConnected() error
	LDs() ([]string, error)
}
