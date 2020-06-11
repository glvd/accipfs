package remote

import (
	"fmt"
	"net"
	"time"

	"github.com/glvd/accipfs/core"
	"github.com/portmapping/go-reuse"
)

type remoteNode struct {
	conn  net.Conn
	addrs []core.Addr
}

// Addrs ...
func (r *remoteNode) Addrs() []core.Addr {
	panic("implement me")
}

// ID ...
func (r *remoteNode) ID() string {
	panic("implement me")
}

// Protocol ...
func (r *remoteNode) Protocol() string {
	panic("implement me")
}

func connectTo(addrs []core.Addr, bindPort int, timeout time.Duration) (net.Conn, error) {
	local := net.TCPAddr{
		IP:   net.IPv4zero,
		Port: bindPort,
	}
	for _, addr := range addrs {
		conn, err := reuse.DialTimeOut("tcp", local.String(), addr.String(), timeout*time.Second)
		if err != nil {
			continue
		}
		return conn, nil
	}
	return nil, fmt.Errorf("all connect failed")
}

func node(conn net.Conn, addrs []core.Addr) core.Node {
	//get info from remote
	return &remoteNode{
		conn:  conn,
		addrs: addrs,
	}
}

// Connect ...
func (r *remoteNode) Connect() (net.Conn, error) {
	connectTo(r.addrs)
}