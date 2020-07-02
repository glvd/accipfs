package core

// NodeManager ...
type NodeManager interface {
	Close()
	Push(n Node)
	Range(f func(key, val []byte) bool)
	HandleConn(c interface{})
	Store() error
	Load() error
}
