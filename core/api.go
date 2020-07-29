package core

import (
	"context"
	"time"
)

// DataStoreReq ...
type DataStoreReq struct {
}

// DataStoreResp ...
type DataStoreResp struct {
	Pins []string
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
	Addrs     []string
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

const (
	// AddNone ...
	AddNone AddType = iota
	// AddOnlyInfo ...
	AddOnlyInfo
	// AddOnlyFile ...
	AddOnlyFile
	// AddBoth ...
	AddBoth
)

// AddReq ...
type AddReq struct {
	//Setting options.UnixfsAddSettings
	Type  AddType
	JSNFO string
	Data  []byte
	Hash  string
}

// AddResp ...
type AddResp struct {
	IsSuccess bool
	Hash      string
}

// UploadReq ...
type UploadReq struct {
	Path string
	//Option options.UnixfsAddOption
}

// UploadResp ...
type UploadResp struct {
	Hash string
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
	ByID    bool
	Names   []string
	Addrs   []string
	Timeout time.Duration
}

// NodeLinkResp ...
type NodeLinkResp struct {
	Err       error
	NodeInfos []NodeInfo
}

// API ...
type API interface {
	Ping(ctx context.Context, req *PingReq) (*PingResp, error)
	ID(ctx context.Context, req *IDReq) (*IDResp, error)
	Add(ctx context.Context, req *AddReq) (*AddResp, error)
	NodeAPI() NodeAPI
	DataStoreAPI() DataStoreAPI
}

// NodeAPI ...
type NodeAPI interface {
	Add(ctx context.Context, req *AddReq) (*AddResp, error)
	Link(ctx context.Context, req *NodeLinkReq) (*NodeLinkResp, error)
	Unlink(ctx context.Context, req *NodeUnlinkReq) (*NodeUnlinkResp, error)
	List(ctx context.Context, req *NodeListReq) (*NodeListResp, error)
	NodeAddrInfo(req *AddrReq) (*AddrResp, error)
}

// DataStoreAPI ...
type DataStoreAPI interface {
	PinLs(ctx context.Context, req *DataStoreReq) (*DataStoreResp, error)
	UploadFile(ctx context.Context, req *UploadReq) (*UploadResp, error)
}
