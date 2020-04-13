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

// UploadReq ...
type UploadReq struct {
}

// UploadResp ...
type UploadResp struct {
}

// BustLinker ...
type BustLinker interface {
	Ping(r *http.Request, req *PingReq, resp *PingResp) error
	Upload(r *http.Request, req *UploadReq, resp *UploadResp) error
}
