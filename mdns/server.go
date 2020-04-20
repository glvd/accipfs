package mdns

import "net"

// Server ...
type Server interface {
}

type server struct {
	cfg  *OptionConfig
	conn []*net.UDPConn
}
