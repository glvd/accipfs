package core

// NodeManager ...
type NodeManager interface {
	Push(node Node)
	Range(f func(key string, node Node) bool)
	Store() error
	Load() error
}
