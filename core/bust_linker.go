package core

import "github.com/libp2p/go-libp2p-core/peer"

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
	AddrInfo  *peer.AddrInfo
	DataStore *DataStoreInfo
	Contract  *ContractInfo
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

// AddTypePeer ...
const AddTypePeer AddType = 0x01

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
