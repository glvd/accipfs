package node

import "github.com/glvd/accipfs/core"

type node struct {
	id       string
	addrs    []core.Addr
	protocol string
}

// Addrs ...
func (n node) Addrs() []core.Addr {
	return n.addrs
}

// ID ...
func (n node) ID() string {
	return n.id
}

// Protocol ...
func (n node) Protocol() string {
	return n.protocol
}
