package node

import (
	"github.com/glvd/accipfs/core"
	"github.com/glvd/accipfs/general"
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

func (n *node) recv() {

}

func (n *node) send() {

}

func nodeRun(node *node) (core.Node, error) {
	go node.running()
	return node, nil
}

// AcceptNode ...
func AcceptNode(conn net.Conn) (core.Node, error) {
	addr := conn.RemoteAddr()
	ip, port := general.SplitIP(addr.String())
	return nodeRun(&node{
		id: "", //todo
		addrs: []core.Addr{
			{
				Protocol: addr.Network(),
				IP:       ip,
				Port:     port,
			},
		},
		conn: conn,
	})
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

func (n *node) running() {

}
