package mdns

import (
	"go.uber.org/atomic"
	"net"
)

// Client ...
type Client interface {
}
type client struct {
	cfg  *OptionConfig
	conn []*net.UDPConn
	stop *atomic.Bool
}
