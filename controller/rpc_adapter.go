package controller

import (
	"github.com/glvd/accipfs/core"
	"net/http"
)

// JSONRPCAdapter ...
type JSONRPCAdapter interface {
	//Ping(r *http.Request, req *core.PingReq, resp *core.PingResp) error
	ID(r *http.Request, req *core.IDReq, resp *core.IDResp) error
	Add(r *http.Request, req *core.AddReq, resp *core.AddResp) error
	Get(r *http.Request, req *core.GetReq, resp *core.GetResp) error
	Pay(r *http.Request, req *core.PayReq, resp *core.PayResp) error
}

type adapter struct {
	api        core.API
	controller *Controller
}

// ID ...
func (a adapter) ID(r *http.Request, req *core.IDReq, resp *core.IDResp) error {
	id, err := a.api.ID(req)
	if err != nil {
		return err
	}
	resp = id
	return nil
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

func newAdapter(api core.API) JSONRPCAdapter {
	return &adapter{
		api: api,
	}
}
