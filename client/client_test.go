package client

import (
	"context"
	httpapi "github.com/ipfs/go-ipfs-http-client"
	ma "github.com/multiformats/go-multiaddr"
	"testing"
)

var _api = defaultAPI()

func defaultAPI() *httpapi.HttpApi {
	multiaddr, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/5001")
	if err != nil {
		panic(err)
	}
	api, err := httpapi.NewApi(multiaddr)
	if err != nil {
		panic(err)
	}
	return api
}

func TestIPFSPinLs(t *testing.T) {
	ls, err := _api.Pin().Ls(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	for ps := range ls {
		t.Logf("%+v",ps.Path().String())
	}

}
