package service

import (
	"github.com/portmapping/go-reuse"
	"net"
)

type tcpListener struct {
	protocol string
	bindPort int
	port     int
	connBack func(conn net.Conn)
}

// NewLinkListener ...
func NewLinkListener(port int, bindPort int, cb func(conn net.Conn)) Listener {
	return &tcpListener{
		protocol: "tcp",
		bindPort: bindPort,
		port:     port,
		connBack: cb,
	}
}

// Listen ...
func (h *tcpListener) Listen() error {
	local := &net.TCPAddr{
		IP:   net.IPv4zero,
		Port: h.port,
	}
	tcp, err := reuse.ListenTCP(h.protocol, local)
	if err != nil {
		return err
	}
	for {
		conn, err := tcp.Accept()
		if err != nil {
			continue
		}
		if h.connBack != nil {
			go h.connBack(conn)
			continue
		}
		//no callback closed
		conn.Close()
	}
}
