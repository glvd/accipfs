package client

import (
	"context"
	"github.com/glvd/accipfs/core"
	files "github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"os"
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
	stat, e := os.Stat(req.Path)
	if e != nil {
		return &core.UploadResp{}, e
	}
	var node files.Node
	//var err error
	if !stat.IsDir() {
		file, e := os.Open(req.Path)
		if e != nil {
			return &core.UploadResp{}, e
		}
		node = files.NewReaderFile(file)
	} else {
		sf, e := files.NewSerialFile(req.Path, false, stat)
		if e != nil {
			return &core.UploadResp{}, e
		}
		node = sf
	}
	if req.Option == nil {
		req.Option = func(settings *options.UnixfsAddSettings) error {
			settings.Pin = true
			return nil
		}
	}

	resolved, e := c.node.Unixfs().Add(context.TODO(), node, req.Option)
	if e != nil {
		return &core.UploadResp{}, e
	}
	return &core.UploadResp{
		Hash: resolved.Cid().String(),
	}, nil
}
