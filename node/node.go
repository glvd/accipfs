package node

import (
	"github.com/glvd/accipfs/core"
	"github.com/portmapping/go-reuse"
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

// ConnectToNode ...
func ConnectToNode(addr core.Addr, bind int) (core.Node, error) {
	tcp, err := reuse.DialTCP(addr.Protocol, &net.TCPAddr{
		IP:   net.IPv4zero,
		Port: bind,
	}, addr.TCP())
	if err != nil {
		return nil, err
	}
	return &node{
		id:    "",
		addrs: []core.Addr{addr},
		conn:  tcp,
	}, nil
}

// AcceptNode ...
func AcceptNode(conn net.Conn) (core.Node, error) {
	addr := conn.RemoteAddr()

	return &node{
		id: "",
		addrs: []core.Addr{
			{addr.Network()},
		},
		conn: tcp,
	}, nil
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
