package core

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
	AddrInfo  *AddrInfo
}

// AddrReq ...
type AddrReq struct {
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
