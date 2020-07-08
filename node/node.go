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

type node struct {
	scdt.Connection
	local    peer.AddrInfo
	remoteID *atomic.String
	remote   peer.AddrInfo
	addrInfo *core.AddrInfo
	api      core.API
}

type jsonNode struct {
	ID    string
	Addrs []ma.Multiaddr
	peer.AddrInfo
}

var _ core.Node = &node{}

// IsClosed ...
func (n *node) IsClosed() bool {
	return n.Connection.IsClosed()
}

// IPFSAddrInfo ...
func (n *node) IPFSAddrInfo() (peer.AddrInfo, error) {
	addrInfo, err := n.addrInfoRequest()
	if err != nil {
		return peer.AddrInfo{}, err
	}
	return addrInfo.IPFSAddrInfo, nil
}

// Marshal ...
func (n *node) Marshal() ([]byte, error) {
	addrInfo, err := n.addrInfoRequest()
	if err != nil {
		return nil, err
	}
	return addrInfo.MarshalJSON()
}

// Unmarshal ...
func (n *node) Unmarshal(bytes []byte) error {
	if n.addrInfo == nil {
		n.addrInfo = new(core.AddrInfo)
	}
	return n.addrInfo.UnmarshalJSON(bytes)
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

// CoreNode ...
func CoreNode(conn net.Conn, api core.API) (core.Node, error) {
	n := &node{
		api:        api,
		Connection: scdt.Connect(conn),
	}
	netAddr, err := mnet.FromNetAddr(conn.RemoteAddr())
	if err != nil {
		return nil, err
	}
	n.AppendAddr(netAddr)
	if err := n.doFirst(); err != nil {
		return nil, err
	}
	return n, nil
}

// ConnectNode ...
func ConnectNode(addr ma.Multiaddr, bind int, api core.API) (core.Node, error) {
	localAddr, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", bind))
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

	n := defaultAPINode(conn, api)
	n.AppendAddr(addr)
	if err := n.doFirst(); err != nil {
		return nil, err
	}
	return n, nil
}

func defaultAPINode(c net.Conn, api core.API) *node {
	conn := scdt.Connect(c, func(c *scdt.Config) {
		c.Timeout = 30 * time.Second
	})
	n := &node{
		api:        api,
		Connection: conn,
	}

	conn.RecvCustomData(func(message *scdt.Message) ([]byte, bool) {
		fmt.Printf("recv data:%+v", message)
		switch message.CustomID {
		case InfoRequest:
			return n.RecvDataRequest(message)
		}
		return nil, false
	})

	return n
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
func (n *node) Info() (core.NodeInfo, error) {
	msg, b := n.Connection.SendCustomDataOnWait(InfoRequest, nil)
	var nodeInfo core.NodeInfo
	fmt.Printf("recved msg:%v,data:%s\n", msg, msg.Data)
	if b && msg.DataLength != 0 {
		err := json.Unmarshal(msg.Data, &nodeInfo)
		if err != nil {
			return nodeInfo, nil
		}
		return nodeInfo, nil
	}
	return nodeInfo, errors.New("data not found")
}

// GetDataRequest ...
func (n *node) GetDataRequest() {

}

// RecvDataRequest ...
func (n *node) RecvDataRequest(message *scdt.Message) ([]byte, bool) {
	addrInfo, err := n.addrInfoRequest()
	if err != nil {
		return nil, true
	}
	nodeInfo := &core.NodeInfo{
		ID:              addrInfo.ID,
		PublicKey:       addrInfo.PublicKey,
		Addrs:           n.Addrs(),
		IPFSAddrInfo:    addrInfo.IPFSAddrInfo,
		AgentVersion:    "", //todo
		ProtocolVersion: "", //todo
	}
	json := nodeInfo.JSON()
	return []byte(json), true
}

func (n *node) addrInfoRequest() (*core.AddrInfo, error) {
	if n.addrInfo != nil {
		return n.addrInfo, nil
	}
	id, err := n.api.NodeAPI().NodeAddrInfo(&core.AddrReq{})
	if err != nil {
		return nil, err
	}
	n.addrInfo = id.AddrInfo
	return n.addrInfo, nil
}

func (n *node) doFirst() error {
	if _, err := n.addrInfoRequest(); err != nil {
		return err
	}
	return nil
}
