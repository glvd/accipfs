package node

import (
	"fmt"
	"net"
	"time"

	"github.com/glvd/accipfs/basis"
	"github.com/glvd/accipfs/core"
	"github.com/godcong/scdt"
	"github.com/portmapping/go-reuse"
)

const maxByteSize = 65520

type jsonNode struct {
	Addrs []core.Addr `json:"addrs"`
}

//temp Data
type nodeLocal struct {
	id *string
}

type node struct {
	scdt.Connection
	api   core.API
	addrs []core.Addr
}

var _ core.Node = &node{}
var heartBeatTimer = 15 * time.Second

// IsClosed ...
func (n *node) IsClosed() bool {
	return n.Connection.IsClosed()
}

// Close ...
func (n *node) Close() (err error) {
	if n.Connection != nil {
		n.Connection.Close()
	}
	return
}

// Verify ...
func (n *node) Verify() bool {
	return true
}

// AcceptNode ...
func AcceptNode(conn net.Conn, api core.API) (core.Node, error) {
	addr := conn.RemoteAddr()
	ip, port := basis.SplitIP(addr.String())
	n := defaultNode(conn)
	n.SetAPI(api)
	n.AppendAddr(core.Addr{
		Protocol: "tcp",
		IP:       ip,
		Port:     port,
	})
	return n, nil
}

// ConnectNode ...
func ConnectNode(addr core.Addr, bind int, api core.API) (core.Node, error) {
	conn, err := reuse.DialTCP(addr.Protocol, &net.TCPAddr{
		IP:   net.IPv4zero,
		Port: bind,
	}, addr.TCP())
	if err != nil {
		return nil, err
	}
	n := defaultNode(conn)
	n.SetAPI(api)
	n.AppendAddr(addr)
	return n, nil
}

func defaultNode(c net.Conn) *node {
	conn := scdt.Connect(c, func(c *scdt.Config) {
		c.Timeout = 30 * time.Second
	})
	return &node{
		api:        nil,
		Connection: conn,
	}
}

// AppendAddr ...
func (n *node) AppendAddr(addrs ...core.Addr) {
	if addrs != nil {
		n.addrs = append(n.addrs, addrs...)
	}
}

// SetAPI ...
func (n *node) SetAPI(api core.API) {
	n.api = api
}

// Addrs ...
func (n node) Addrs() []core.Addr {
	return n.addrs
}

// ID ...
func (n *node) ID() string {
	id, err := n.Connection.RemoteID()
	if err != nil {
		return ""
	}
	return id
}

// Info ...
func (n *node) Info() core.NodeInfo {
	msg, b := n.Connection.SendCustomDataOnWait(0x01, []byte("get info from remote"))

	if b {
		fmt.Println("result msg", string(msg.Data))
		return core.NodeInfo{
			ID:   "recv id",
			Type: core.NodeAccelerate,
		}
	}
	return core.NodeInfo{}
}
