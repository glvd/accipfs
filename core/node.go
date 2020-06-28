package core

// RecvCBFunc ...
type RecvCBFunc func(id string, v interface{}) ([]byte, error)

// Node ...
type Node interface {
	ID() string
	Addrs() []Addr
	Info() NodeInfo
	Close() (err error)
	IsClosed() bool
}
