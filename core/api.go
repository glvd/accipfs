package core

// API ...
type API interface {
	Ping(req *PingReq) (*PingResp, error)
	ID(req *IDReq) (*IDResp, error)
	Add(req *AddReq) (*AddResp, error)
	NodeAPI() NodeAPI
}

// NodeListReq ...
type NodeListReq struct {
}

// NodeListResp ...
type NodeListResp struct {
}

// NodeUnlinkReq ...
type NodeUnlinkReq struct {
}

// NodeUnlinkResp ...
type NodeUnlinkResp struct {
}

// NodeAPI ...
type NodeAPI interface {
	Link(req *NodeLinkReq) (*NodeLinkResp, error)
	Unlink(req *NodeUnlinkReq) (*NodeUnlinkResp, error)
	List(req *NodeListReq) (*NodeListResp, error)
	NodeAddrInfo(req *AddrReq) (*AddrResp, error)
}
