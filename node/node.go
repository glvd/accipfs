package node

import (
	"fmt"
	"net"
	"time"

	"github.com/glvd/accipfs/core"
	"github.com/godcong/scdt"
	ma "github.com/multiformats/go-multiaddr"
	mnet "github.com/multiformats/go-multiaddr-net"
)

const (
	// InfoRequest ...
	InfoRequest = iota + 1
)

type jsonNode struct {
	Addrs []string `json:"addrs"`
}

//temp Data
type nodeLocal struct {
	id *string
}

type node struct {
	scdt.Connection
	api   core.API
	addrs []ma.Multiaddr
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
	n := defaultNode(conn)
	n.SetAPI(api)
	netAddr, err := mnet.FromNetAddr(conn.RemoteAddr())
	if err != nil {
		return nil, err
	}
	n.AppendAddr(netAddr)
	return n, nil
}

// ConnectNode ...
func ConnectNode(addr ma.Multiaddr, bind int, api core.API) (core.Node, error) {

	conn, err := mnet.Dial(addr)
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
func (n *node) AppendAddr(addrs ...ma.Multiaddr) {
	if addrs != nil {
		n.addrs = append(n.addrs, addrs...)
	}
}

// SetAPI ...
func (n *node) SetAPI(api core.API) {
	n.api = api
}

// Addrs ...
func (n node) Addrs() []ma.Multiaddr {
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
	msg, b := n.Connection.SendOnWait([]byte("get info from remote"))

	if b {
		fmt.Println("result msg", string(msg.Data))
		return core.NodeInfo{
			ID:   "recv id",
			Type: core.NodeAccelerate,
		}
	}
	return core.NodeInfo{}
}

// GetDataRequest ...
func (n *node) GetDataRequest() {

}

// RecvDataRequest ...
func (n *node) RecvDataRequest(id uint16, cb scdt.RecvCallbackFunc) {
	//todo
}
