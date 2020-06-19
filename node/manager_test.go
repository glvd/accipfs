package node

import (
	"github.com/glvd/accipfs/basis"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	alog "github.com/glvd/accipfs/log"
	"net"
	"testing"
	"time"
)

var testConfig = config.Default()

func init() {
	alog.InitLog()
	testConfig.Path = ""
	testConfig.Node.BackupSeconds = 3 * time.Second
}

func TestNew(t *testing.T) {
	m := New(testConfig)
	err := m.Load()
	if err != nil {
		t.Log(err)
	}

	for i := 0; i < 100; i++ {
		m.Push(&node{
			id: basis.UUID(),
			addrs: []core.Addr{
				{
					Protocol: "tcp",
					IP:       net.IPv4zero,
					Port:     16005,
				},
			},
			conn: nil,
		})
	}
	time.Sleep(5 * time.Second)
}
