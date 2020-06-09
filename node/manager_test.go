package node

import (
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/glvd/accipfs/general"
	alog "github.com/glvd/accipfs/log"
	"net"
	"testing"
	"time"
)

var testConfig = config.Default()

func init() {
	alog.InitLog()
	testConfig.Path = ""
	testConfig.Node.BackupSeconds = 5 * 60 * time.Second
}

func TestNew(t *testing.T) {
	m := New(testConfig)
	err := m.Load()
	if err != nil {
		t.Log(err)
	}

	m.Push(&node{
		id: general.UUID(),
		addrs: []core.Addr{
			{
				Protocol: "tcp",
				IP:       net.IPv4zero,
				Port:     16005,
			},
		},
		conn: nil,
	})

	m.Store()
}
