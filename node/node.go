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
	local          core.SafeLocalData
	remoteID       *atomic.String
	remote         peer.AddrInfo
	remoteNodeInfo *core.NodeInfo
	//addrInfo       *core.AddrInfo
	//api            core.API
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

// DataStoreInfo ...
func (n *node) DataStoreInfo() (core.DataStoreInfo, error) {
	addrInfo, err := n.addrInfoRequest()
	if err != nil {
		return core.DataStoreInfo{}, err
	}
	return addrInfo.DataStore, nil
}

// Marshal ...
func (n *node) Marshal() ([]byte, error) {
	return n.local.Marshal()
}

// Unmarshal ...
func (n *node) Unmarshal(bytes []byte) error {
	return n.local.Unmarshal(bytes)
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
func CoreNode(conn net.Conn, local core.SafeLocalData) (core.Node, error) {
	n := defaultAPINode(conn, local, 30*time.Second)
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
func ConnectNode(addr ma.Multiaddr, bind int, local core.SafeLocalData) (core.Node, error) {
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
	n := defaultAPINode(conn, local, 0)
	n.AppendAddr(addr)
	if err := n.doFirst(); err != nil {
		return nil, err
	}
	return n, nil
}

func defaultAPINode(c net.Conn, local core.SafeLocalData, duration time.Duration) *node {
	conn := scdt.Connect(c, func(c *scdt.Config) {
		c.Timeout = duration
		c.CustomIDer = func() string {
			return local.Data().Node.ID
		}
	})
	n := &node{
		local:      local,
		Connection: conn,
	}

	conn.Recv(func(message *scdt.Message) ([]byte, bool, error) {
		fmt.Printf("recv data:%+v", message)
		return []byte("recv called"), true, errors.New("not data")
	})
	conn.RecvCustomData(func(message *scdt.Message) ([]byte, bool, error) {
		fmt.Printf("recv custom data:%+v", message)
		switch message.CustomID {
		case InfoRequest:
			request, b, err := n.RecvDataRequest(message)
			return request, b, err
		}
		return []byte("recv custom called"), true, errors.New("wrong case")
	})

	return n
}

// AppendAddr ...
func (n *node) AppendAddr(addrs ...ma.Multiaddr) {
	if addrs != nil {
		n.remote.Addrs = append(n.remote.Addrs, addrs...)
	}
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
	if n.remoteNodeInfo != nil {
		return *n.remoteNodeInfo, nil
	}
	return n.GetInfoDataRequest()
}

// GetDataRequest ...
func (n *node) GetInfoDataRequest() (core.NodeInfo, error) {
	msg, b := n.Connection.SendCustomDataOnWait(InfoRequest, nil)
	var nodeInfo core.NodeInfo
	fmt.Printf("recved msg:%v\n", msg)
	if b && msg.DataLength != 0 {
		fmt.Printf("msg data:%v\n", string(msg.Data))
		err := json.Unmarshal(msg.Data, &nodeInfo)
		if err != nil {
			return nodeInfo, err
		}
		return nodeInfo, nil
	}
	n.remoteNodeInfo = &nodeInfo
	return nodeInfo, errors.New("data not found")
}

// RecvIndexSyncRequest ...
func (n *node) RecvIndexSyncRequest() ([]byte, bool, error) {
	panic("//todo")
}

// RecvNodeListRequest ...
func (n *node) RecvNodeListRequest() ([]byte, bool, error) {
	panic("//todo")
}

// RecvDataRequest ...
func (n *node) RecvDataRequest(message *scdt.Message) ([]byte, bool, error) {
	fmt.Printf("request %v\n", message)
	addrInfo, err := n.addrInfoRequest()
	if err != nil {
		return nil, true, err
	}
	nodeInfo := &core.NodeInfo{
		AddrInfo:        *addrInfo,
		AgentVersion:    "", //todo
		ProtocolVersion: "", //todo
	}
	json := nodeInfo.JSON()
	log.Infow("node info", "json", json)
	return []byte(json), true, nil
}

func (n *node) addrInfoRequest() (*core.AddrInfo, error) {
	data := n.local.Data()
	return &data.Node.AddrInfo, nil
	//if n.addrInfo != nil {
	//	return n.addrInfo, nil
	//}
	//id, err := n.api.NodeAPI().NodeAddrInfo(&core.AddrReq{})
	//if err != nil {
	//	return nil, err
	//}
	//n.addrInfo = id.AddrInfo
	//return &n.local.Node.AddrInfo
}

func (n *node) doFirst() error {
	if _, err := n.addrInfoRequest(); err != nil {
		return err
	}
	return nil
}
