package mdns

import "net"

// Client ...
type Client interface {
}
type client struct {
	conn []*net.UDPConn
}
