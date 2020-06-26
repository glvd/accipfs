package core

type RecvCBFunc func(id string, v interface{}) ([]byte, error)

// Node ...
type Node interface {
	ID() string
	RecvCallback(fn RecvCBFunc)
	Addrs() []Addr
	Info() NodeInfo
	Ping() error
	IsConnecting() bool
	Closed(f func(Node))
	Close() (err error)
	IsClosed() bool
}
