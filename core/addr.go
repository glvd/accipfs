package core

import (
	"net"
	"strconv"
)

// Addr ...
type Addr struct {
	Protocol string
	IP       net.IP
	Port     int
}

// Network ...
func (addr Addr) Network() string {
	return addr.Protocol
}

// String ...
func (addr Addr) String() string {
	return net.JoinHostPort(addr.IP.String(), strconv.Itoa(addr.Port))
}

// UDP ...
func (addr Addr) UDP() *net.UDPAddr {
	return &net.UDPAddr{
		IP:   addr.IP,
		Port: addr.Port,
	}
}

// TCP ...
func (addr Addr) TCP() *net.TCPAddr {
	return &net.TCPAddr{
		IP:   addr.IP,
		Port: addr.Port,
	}
}

// IsZero ...
func (addr Addr) IsZero() bool {
	return addr.Protocol == "" && addr.IP.Equal(net.IPv4zero) && addr.Port == 0
}
