package service

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/rpc"
	"testing"
)

func TestNodeClientETH_NewAccount(t *testing.T) {
	var inf interface{}
	cancelCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := rpc.DialContext(cancelCtx, "http://127.0.0.1:8545")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	err = client.Call(&inf, "admin_addPeer", "enode://cd8bb3420de832d2d48e8c5ca70d83cac6a2d01fde8f2259cb725ad9d92d2bd52200d817c95e6407c09b8660806132c40f1e6dab3f04411407144ec5d3c1060c@10.250.52.62:30303?discport=0")
	if err != nil {
		t.Fatal(err)
	}

	marshal, err := json.Marshal(inf)
	if err != nil {
		return
	}
	t.Logf("info:%+v", string(marshal))

}
