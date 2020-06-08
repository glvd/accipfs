package remote

import (
	"net"

	"github.com/glvd/accipfs/core"
)

type remoteConnect struct {
	conn net.Conn
}

func ConnectTo(addrs []core.Addr) net.Conn {
	for _, addr := range addrs {
		reuse.Dial()
	}
}
