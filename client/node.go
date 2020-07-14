package client

import "github.com/glvd/accipfs/core"

// NodeAPI ...
func (c *client) NodeAPI() core.NodeAPI {
	return c
}

// Unlink ...
func (c *client) Unlink(req *core.NodeUnlinkReq) (resp *core.NodeUnlinkResp, err error) {
	resp = new(core.NodeUnlinkResp)
	err = c.doPost("node/unlink", req, resp)
	return
}

// NodeList ...
func (c *client) List(req *core.NodeListReq) (resp *core.NodeListResp, err error) {
	resp = new(core.NodeListResp)
	err = c.doPost("node/list", req, resp)
	return
}

// NodeList ...
func NodeList(req *core.NodeListReq) (resp *core.NodeListResp, err error) {
	return DefaultClient.NodeAPI().List(req)
}

// Link ...
func (c *client) Link(req *core.NodeLinkReq) (resp *core.NodeLinkResp, err error) {
	resp = new(core.NodeLinkResp)
	err = c.doPost("/node/link", req, resp)
	return
}

// NodeLink ...
func NodeLink(req *core.NodeLinkReq) (resp *core.NodeLinkResp, err error) {
	return DefaultClient.NodeAPI().Link(req)
}
