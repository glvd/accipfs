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

// Verify ...
func (n *node) Verify() bool {
	return true
}

var _ core.Node = &node{}

// Close ...
func (n *node) Close() (err error) {
	if n.conn != nil {
		err = n.conn.Close()
		n.conn = nil
	}
	return
}

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
