package mdns

import "net"

// Server ...
type Server interface {
}

type server struct {
	conn []*net.UDPConn
}
