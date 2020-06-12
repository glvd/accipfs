package core

// Node ...
type Node interface {
	Addrs() []Addr
	ID() string
	Verify() bool
	Info() NodeInfo
}

// LocalNode ...
type LocalNode interface {
	Addr() Addr
}
