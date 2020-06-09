package node

import (
	"fmt"
	"github.com/glvd/accipfs/core"
	"net"
)

type jsonNode struct {
	Addrs []core.Addr `json:"addrs"`
}

type node struct {
	id    string
	addrs []core.Addr
	conn  net.Conn
}

var _ core.Node = &node{}

// Connect ...
func (n *node) Connect() (net.Conn, error) {
	if n.conn != nil {
		return n.conn, nil
	}
	//todo
	return nil, fmt.Errorf("filed to connect")
}

// Addrs ...
func (n node) Addrs() []core.Addr {
	return n.addrs
}

// ID ...
func (n node) ID() string {
	return n.id
}
