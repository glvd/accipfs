package core

// Node ...
type Node interface {
	Addrs() []Addr
	Info() NodeInfo
	Ping() error
}
