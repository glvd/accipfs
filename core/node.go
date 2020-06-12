package core

// Node ...
type Node interface {
	Addrs() []Addr
	ID() string
	Verify() bool
	//Connect() (net.Conn, error)
	//Close() error
}
