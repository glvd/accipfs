package client

import (
	"context"
	"github.com/glvd/accipfs/core"
)

// NodeAPI ...
func (c *client) NodeAPI() core.NodeAPI {
	return c
}

// Unlink ...
func (c *client) Unlink(ctx context.Context, req *core.NodeUnlinkReq) (resp *core.NodeUnlinkResp, err error) {
	resp = new(core.NodeUnlinkResp)
	err = c.doPost(ctx, "node/unlink", req, resp)
	return
}

// NodeList ...
func (c *client) List(ctx context.Context, req *core.NodeListReq) (resp *core.NodeListResp, err error) {
	resp = new(core.NodeListResp)
	err = c.doPost(ctx, "node/list", req, resp)
	return
}

// NodeList ...
func NodeList(ctx context.Context, req *core.NodeListReq) (resp *core.NodeListResp, err error) {
	return DefaultClient.NodeAPI().List(ctx, req)
}

// Link ...
func (c *client) Link(ctx context.Context, req *core.NodeLinkReq) (resp *core.NodeLinkResp, err error) {
	resp = new(core.NodeLinkResp)
	err = c.doPost(ctx, "/node/link", req, resp)
	return
}

// NodeLink ...
func NodeLink(ctx context.Context, req *core.NodeLinkReq) (resp *core.NodeLinkResp, err error) {
	return DefaultClient.NodeAPI().Link(ctx, req)
}
