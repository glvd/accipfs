package core

// Node ...
type Node interface {
	Addrs() []Addr
	ID() string
	Protocol() string
}
