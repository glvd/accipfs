package service

import "net"

// Listener ...
type Listener interface {
	Listen() error
	Accept(func(conn net.Conn))
}
