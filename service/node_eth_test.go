package service

import (
	"context"
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
	err = client.Call(&inf, "personal_newAccount")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("info:%+v", inf)

}
