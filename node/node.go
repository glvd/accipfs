package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/glvd/accipfs/core"
	"github.com/godcong/scdt"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	mnet "github.com/multiformats/go-multiaddr-net"
	"go.uber.org/atomic"
	"net"
	"time"
)

const (
	// AlreadyConnectedRequest ...
	AlreadyConnectedRequest = iota + 1
	// InfoRequest ...
	InfoRequest
	// LDsRequest ...
	LDsRequest
	// PeerRequest ...
	PeerRequest
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

// ErrNoData ...
var ErrNoData = errors.New("no data respond")

// SendClose ...
func (n *node) SendClose() {
	n.Connection.SendClose([]byte("connected"))
}

// IsClosed ...
func (n *node) IsClosed() bool {
	return n.Connection.IsClosed()
}

// Peers ...
func (n *node) Peers() ([]string, error) {
	msg, b := n.Connection.SendCustomDataOnWait(PeerRequest, nil)
	var s []string
	if b {
		if msg.DataLength > 0 {
			err := json.Unmarshal(msg.Data, &s)
			if err != nil {
				return nil, err
			}
			return s, nil
		}
	}
	return nil, ErrNoData
}

// SendPeerRequest ...
func (n *node) SendPeerRequest() ([]core.NodeInfo, error) {
	msg, b := n.Connection.SendCustomDataOnWait(PeerRequest, nil)
	var s []core.NodeInfo
	if b {
		if msg.DataLength > 0 {
			err := json.Unmarshal(msg.Data, &s)
			if err != nil {
				return nil, err
			}
			return s, nil
		}
	}
	return nil, ErrNoData
}

// LDs ...
func (n *node) LDs() ([]string, error) {
	msg, b := n.Connection.SendCustomDataOnWait(LDsRequest, nil)
	var s []string
	if b {
		if msg.DataLength > 0 {
			fmt.Println("recv lds", string(msg.Data))
			err := json.Unmarshal(msg.Data, &s)
			if err != nil {
				return nil, err
			}
			return s, nil
		}
	}
	return nil, ErrNoData
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
		//fmt.Printf("recv custom data:%+v\n", message)
		switch message.CustomID {
		case InfoRequest:
			request, b, err := n.RecvInfoRequest(message)
			return request, b, err
		case PeerRequest:
			request, b, err := n.RecvPeerRequest(message)
			return request, b, err
		case LDsRequest:
			request, b, err := n.RecvLDsRequest(message)
			return request, b, err
		case AlreadyConnectedRequest:
			conn.Close()
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

// SendConnected ...
func (n *node) SendConnected() error {
	n.SendCustomData(AlreadyConnectedRequest, []byte("connected"))
	return nil
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

// GetInfo ...
func (n *node) GetInfo() (core.NodeInfo, error) {
	if n.remoteNodeInfo != nil {
		return *n.remoteNodeInfo, nil
	}
	return n.SendInfoRequest()
}

// GetDataRequest ...
func (n *node) SendInfoRequest() (core.NodeInfo, error) {
	msg, b := n.Connection.SendCustomDataOnWait(InfoRequest, nil)
	var nodeInfo core.NodeInfo
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
	addrs := n.local.Data().Addrs
	marshal, err := json.Marshal(addrs)
	if err != nil {
		return nil, false, err
	}
	return marshal, true, nil
}

// RecvInfoRequest ...
func (n *node) RecvInfoRequest(message *scdt.Message) ([]byte, bool, error) {
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
}

func (n *node) doFirst() error {
	for i := 0; i < 3; i++ {
		if !n.local.Data().Initialized {
			time.Sleep(3 * time.Second)
			continue
		}
		return nil
	}
	return errors.New("first init error")
}

// RecvPeerRequest ...
func (n *node) RecvPeerRequest(message *scdt.Message) ([]byte, bool, error) {
	peers := n.local.Data().Nodes
	marshal, err := json.Marshal(peers)
	if err != nil {
		return nil, false, err
	}
	return marshal, true, nil
}

// RecvLDsRequest ...
func (n *node) RecvLDsRequest(message *scdt.Message) ([]byte, bool, error) {
	lds := n.local.Data().LDs
	var ret []string
	for ld := range lds {
		ret = append(ret, ld)
	}
	marshal, err := json.Marshal(ret)
	if err != nil {
		return nil, false, err
	}
	return marshal, true, nil
}
