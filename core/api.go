package core

// API ...
type API interface {
	Ping(req *PingReq) (*PingResp, error)
}
