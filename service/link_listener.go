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
func NewLinkListener(protocol string, port int, bindPort int) Listener {
	return &tcpListener{
		protocol: protocol,
		bindPort: bindPort,
		port:     port,
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

// Accept ...
func (h *tcpListener) Accept(f func(conn net.Conn)) {
	if f != nil {
		h.connBack = f
	}
}
