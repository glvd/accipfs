package core

import "net"

// Addr ...
type Addr struct {
	Protocol string
	IP       net.IP
	Port     int
}
