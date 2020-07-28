package client

import (
	"github.com/glvd/accipfs/core"
)

// DataStoreAPI ...
func (c *client) DataStoreAPI() core.DataStoreAPI {
	return c
}

// DataStorePinLs ...
func DataStorePinLs(req *core.DataStoreReq) (resp *core.DataStoreResp, err error) {
	return DefaultClient.DataStoreAPI().PinLs(req)
}

// PinLs ...
func (c *client) PinLs(req *core.DataStoreReq) (resp *core.DataStoreResp, err error) {
	resp = new(core.DataStoreResp)
	err = c.doPost("ds/pin/ls", req, resp)
	return
}

// UploadFile ...
func (c *client) UploadFile(req *core.UploadReq) (resp *core.UploadResp, err error) {
	resp = new(core.UploadResp)
	err = c.doPost("ds/upload", req, resp)
	return
}
