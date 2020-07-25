package client

import (
	httpapi "github.com/ipfs/go-ipfs-http-client"
	ma "github.com/multiformats/go-multiaddr"
	"testing"
)

var api = defaultAPI()

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

func TestPIN(t *testing.T) {
	t.Run("TestPinAdd", TestPinAdd)
	t.Run("TestPinSimple", TestPinSimple)
	t.Run("TestPinRecursive", TestPinRecursive)
	t.Run("TestPinLsIndirect", TestPinLsIndirect)
	t.Run("TestPinLsPrecedence", TestPinLsPrecedence)
	t.Run("TestPinIsPinned", TestPinIsPinned)
}
