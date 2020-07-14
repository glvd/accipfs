package core

import ma "github.com/multiformats/go-multiaddr"

// DataStoreReq ...
type DataStoreReq struct {
}

// DataStoreResp ...
type DataStoreResp struct {
}

// PingReq ...
type PingReq struct {
}

// PingResp ...
type PingResp struct {
	Data string
}

// PayReq ...
type PayReq struct {
}

// PayResp ...
type PayResp struct {
}

// IDReq ...
type IDReq struct {
}

// IDResp ...
type IDResp struct {
	ID        string
	PublicKey string
	Addrs     []ma.Multiaddr
	DataStore DataStoreInfo
}

// AddrReq ...
type AddrReq struct {
	ID string
}

// AddrResp ...
type AddrResp struct {
	AddrInfo AddrInfo
}

// NodeListReq ...
type NodeListReq struct {
}

// NodeListResp ...
type NodeListResp struct {
	Nodes map[string]NodeInfo
}

// NodeUnlinkReq ...
type NodeUnlinkReq struct {
	Peers []string
}

// NodeUnlinkResp ...
type NodeUnlinkResp struct {
}

// ConnectToReq ...
type ConnectToReq struct {
}

// ConnectToResp ...
type ConnectToResp struct {
	Node
}

// AddType ...
type AddType int

// AddReq ...
type AddReq struct {
	JSNFO string
	Data  []byte
}

// AddResp ...
type AddResp struct {
	IsSuccess bool
}

// GetReq ...
type GetReq struct {
}

// GetResp ...
type GetResp struct {
}

// RequestTag ...
type RequestTag int

// NodeLinkReq ...
type NodeLinkReq struct {
	Addrs []string
}

// NodeLinkResp ...
type NodeLinkResp struct {
	Err error
	NodeInfo
}

// API ...
type API interface {
	Ping(req *PingReq) (*PingResp, error)
	ID(req *IDReq) (*IDResp, error)
	Add(req *AddReq) (*AddResp, error)
	NodeAPI() NodeAPI
	DataStoreAPI() DataStoreAPI
}

// NodeAPI ...
type NodeAPI interface {
	Link(req *NodeLinkReq) (*NodeLinkResp, error)
	Unlink(req *NodeUnlinkReq) (*NodeUnlinkResp, error)
	List(req *NodeListReq) (*NodeListResp, error)
	NodeAddrInfo(req *AddrReq) (*AddrResp, error)
}

// DataStoreAPI ...
type DataStoreAPI interface {
	Pins(req *DataStoreReq) (*DataStoreResp, error)
}
