package core

import ma "github.com/multiformats/go-multiaddr"

// AddTypePeer ...
const AddTypePeer AddType = 0x01

// RequestTagLink ...
const RequestTagLink RequestTag = iota

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
	Name      string
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
	AddrInfo *AddrInfo
}

// ConnectedReq ...
type ConnectedReq struct {
	Node
}

// ConnectedResp ...
type ConnectedResp struct {
	Node
}

// ConnectToReq ...
type ConnectToReq struct {
	Addr string
}

// ConnectToResp ...
type ConnectToResp struct {
	Node
}

// AddType ...
type AddType int

// AddReq ...
type AddReq struct {
	AddType
	Node
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

// LinkReq ...
type LinkReq struct {
	Addrs []string
}

// LinkResp ...
type LinkResp struct {
	Addr string
}
