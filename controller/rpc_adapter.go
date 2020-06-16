package controller

import (
	"github.com/glvd/accipfs/core"
	"net/http"
)

type adapter struct {
	api        core.API
	controller *Controller
}

// Ping ...
func (a adapter) Ping(r *http.Request, req *core.PingReq, resp *core.PingResp) error {
	panic("implement me")
}

// ID ...
func (a adapter) ID(r *http.Request, req *core.IDReq, resp *core.IDResp) error {
	panic("implement me")
}

// Add ...
func (a adapter) Add(r *http.Request, req *core.AddReq, resp *core.AddResp) error {
	panic("implement me")
}

// Get ...
func (a adapter) Get(r *http.Request, req *core.GetReq, resp *core.GetResp) error {
	panic("implement me")
}

// Pay ...
func (a adapter) Pay(r *http.Request, req *core.PayReq, resp *core.PayResp) error {
	panic("implement me")
}

func newAdapter(api core.API) core.LinkerRPC {
	return &adapter{
		api: api,
	}
}
