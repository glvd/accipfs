package core

// Node ...
type Node interface {
	ID() string
	Addrs() []Addr
	Info() NodeInfo
	Ping() error
	IsConnecting() bool
}
