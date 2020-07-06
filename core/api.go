package core

// API ...
type API interface {
	Ping(req *PingReq) (*PingResp, error)
	ID(req *IDReq) (*IDResp, error)
	NodeAddrInfo(req *AddrReq) (*AddrResp, error)
	Link(req *LinkReq) (*LinkResp, error)
}
