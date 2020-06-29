package core

import ma "github.com/multiformats/go-multiaddr"

// RecvCBFunc ...
type RecvCBFunc func(id string, v interface{}) ([]byte, error)

// Node ...
type Node interface {
	ID() string
	Addrs() []ma.Multiaddr
	Info() NodeInfo
	Close() (err error)
	IsClosed() bool
}
