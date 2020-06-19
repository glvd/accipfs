package core

// NodeManager ...
type NodeManager interface {
	Push(n Node)
	Range(f func(key string, n Node) bool)
	HandleConn(c interface{})
	Store() error
	Load() error
}
