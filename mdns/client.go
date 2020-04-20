package mdns

import (
	"go.uber.org/atomic"
	"net"
)

// Client ...
type Client interface {
}
type client struct {
	cfg     *OptionConfig
	uniConn []*net.UDPConn
	conn    []*net.UDPConn
	stop    *atomic.Bool
}
