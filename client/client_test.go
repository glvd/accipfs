package client

import (
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"testing"
)

func init() {
	DefaultClient = New(&config.APIConfig{
		Port:        10808,
		Version:     "",
		UseTLS:      false,
		TLS:         config.TLSCertificate{},
		Certificate: nil,
		Timeout:     30,
	})
}

func TestDataStorePinLs(t *testing.T) {
	ls, err := DefaultClient.DataStoreAPI().PinLs(&core.DataStoreReq{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ls)
}
