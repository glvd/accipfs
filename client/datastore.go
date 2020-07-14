package client

import "github.com/glvd/accipfs/core"

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
	err = c.doPost("datastore/pin/ls", req, resp)
	return
}
