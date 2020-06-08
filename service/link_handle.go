package service

import (
	"github.com/portmapping/go-reuse"
	"net"
)

type handleTCP struct {
	protocol string
	port     int
}

// Listen ...
func (h *handleTCP) Listen() {
	local := &net.TCPAddr{
		IP:   net.IPv4zero,
		Port: h.port,
	}
	tcp, err := reuse.ListenTCP(h.protocol, local)
	if err != nil {
		return
	}
	for {
		conn, err := tcp.Accept()
		if err != nil {
			continue
		}

	}
}
