package core

// NodeManager ...
type NodeManager interface {
	Close()
	Push(n Node)
	Range(f func(key string, node Node) bool)
	HandleConn(c interface{})
	Store() error
	Load() error
}
