package client

import (
	"context"
	"github.com/glvd/accipfs/core"
)

// DataStoreAPI ...
func (c *client) DataStoreAPI() core.DataStoreAPI {
	return c
}

// DataStorePinLs ...
func DataStorePinLs(ctx context.Context, req *core.DataStoreReq) (resp *core.DataStoreResp, err error) {
	return DefaultClient.DataStoreAPI().PinLs(ctx, req)
}

// PinLs ...
func (c *client) PinLs(ctx context.Context, req *core.DataStoreReq) (resp *core.DataStoreResp, err error) {
	resp = new(core.DataStoreResp)
	err = c.doPost(ctx, "ds/pin/ls", req, resp)
	return
}

// UploadFile ...
func (c *client) UploadFile(ctx context.Context, req *core.UploadReq) (resp *core.UploadResp, err error) {
	resp = new(core.UploadResp)
	err = c.doPost(ctx, "ds/upload", req, resp)
	return
}
