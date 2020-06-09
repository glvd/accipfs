package core

import "net"

// Node ...
type Node interface {
	Addrs() []Addr
	ID() string
	Connect() (net.Conn, error)
	Close() error
}
