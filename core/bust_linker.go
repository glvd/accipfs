package core

import (
	"net/http"
)

// PingReq ...
type PingReq struct {
}

// PingResp ...
type PingResp struct {
	Resp string
}

// SourceReq ...
type SourceReq struct {
}

// SourceResp ...
type SourceResp struct {
}

// MessageReq ...
type MessageReq struct {
}

// MessageResp ...
type MessageResp struct {
}

// PayReq ...
type PayReq struct {
}

// PayResp ...
type PayResp struct {
}

// BustLinker ...
type BustLinker interface {
	Ping(r *http.Request, req *PingReq, resp *PingResp) error
	Source(r *http.Request, req *SourceReq, resp *SourceResp) error
	Message(r *http.Request, req *MessageReq, resp *MessageResp) error
	Pay(r *http.Request, req *PayReq, resp *PayResp) error
}
