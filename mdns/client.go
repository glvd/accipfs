package mdns

import "net"

// Client ...
type Client interface {
}
type client struct {
	cfg  *OptionConfig
	conn []*net.UDPConn
}
