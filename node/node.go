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

// Close ...
func (n *node) Close() (err error) {
	if n.conn != nil {
		err = n.conn.Close()
		n.conn = nil
	}
	return
}

// Verify ...
func (n *node) Verify() bool {
	return true
}

// Connect ...
func (n *node) ConnectTo() (net.Conn, error) {
	if n.conn != nil {
		return n.conn, nil
	}
	//todo
	return nil, fmt.Errorf("filed to connect")
}

// ConnectBy ...
func (n *node) ConnectBy() {

}

// Addrs ...
func (n node) Addrs() []core.Addr {
	return n.addrs
}

// ID ...
func (n node) ID() string {
	return n.id
}

// Info ...
func (n *node) Info() core.NodeInfo {
	panic("implement me")
}

// Ping ...
func (n *node) Ping() error {
	panic("implement me")
}
