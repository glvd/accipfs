package core

// NodeManager ...
type NodeManager interface {
	Push(node Node)
	Range(f func(node Node))
	Store() error
	Load() error
}
