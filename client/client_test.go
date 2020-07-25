package client

import (
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"testing"
)

func init() {
	DefaultClient = New(config.Default())
}

func TestDataStorePinLs(t *testing.T) {
	ls, err := DefaultClient.DataStoreAPI().PinLs(&core.DataStoreReq{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ls)
}
