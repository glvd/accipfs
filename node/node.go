package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	"go.uber.org/atomic"
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
	local    *peer.AddrInfo
	remoteID *atomic.String
	remote   *peer.AddrInfo
	api      core.API
}

var _ core.Node = &node{}

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
	localAddr, err := ma.NewMultiaddr(fmt.Sprintf("/tcp/%d", bind))
	if err != nil {
		return nil, err
	}
	d := mnet.Dialer{
		Dialer:    net.Dialer{},
		LocalAddr: localAddr,
	}
	conn, err := d.Dial(addr)
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
		n.remote.Addrs = append(n.remote.Addrs, addrs...)
	}
}

// SetAPI ...
func (n *node) SetAPI(api core.API) {
	n.api = api
}

// Addrs ...
func (n node) Addrs() []ma.Multiaddr {
	return n.remote.Addrs
}

// ID ...
func (n *node) ID() string {
	if n.remoteID != nil {
		return n.remoteID.Load()
	}
	id, err := n.Connection.RemoteID()
	if err != nil {
		return ""
	}
	n.remoteID = atomic.NewString(id)
	return id
}

// Info ...
func (n *node) Info() (peer.AddrInfo, error) {
	msg, b := n.Connection.SendCustomDataOnWait(InfoRequest, []byte("get info from remote"))
	var ai peer.AddrInfo
	if b {
		if msg.DataLength != 0 {
			err := json.Unmarshal(msg.Data, &ai)
			if err != nil {
				return ai, nil
			}
			return ai, nil
		}
	}
	return ai, errors.New("data not found")
}

// GetDataRequest ...
func (n *node) GetDataRequest() {

}

// RecvDataRequest ...
func (n *node) RecvDataRequest(id uint16, cb scdt.RecvCallbackFunc) {
	switch id {
	case InfoRequest:

	}
}

func (n *node) infoRequest() peer.AddrInfo {
	id, err := n.api.ID(&core.IDReq{})
	if err != nil {
		return peer.AddrInfo{}
	}
	return *id.AddrInfo
}
