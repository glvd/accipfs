package core

import (
	"net/http"
)

// PingReq ...
type PingReq struct {
	NodeAddress
}

// PingResp ...
type PingResp struct {
	Resp string
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
	NodeInfo
}

// ConnectReq ...
type ConnectReq struct {
	NodeInfo
}

// ConnectResp ...
type ConnectResp struct {
	NodeInfo
}

// ConnectToReq ...
type ConnectToReq struct {
	Addr string
}

// ConnectToResp ...
type ConnectToResp struct {
	NodeInfo
}

// AddType ...
type AddType int

// AddTypePeer ...
const AddTypePeer AddType = 0x01

// AddReq ...
type AddReq struct {
	AddType
	NodeInfo
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

// Linker ...
type Linker interface {
	Ping(r *http.Request, req *PingReq, resp *PingResp) error
	ID(r *http.Request, req *IDReq, resp *IDResp) error
	Connected(r *http.Request, req *ConnectReq, resp *ConnectResp) error
	ConnectTo(r *http.Request, req *ConnectToReq, resp *ConnectToResp) error
	Add(r *http.Request, req *AddReq, resp *AddResp) error
	Get(r *http.Request, req *GetReq, resp *GetResp) error
	Pay(r *http.Request, req *PayReq, resp *PayResp) error
}
