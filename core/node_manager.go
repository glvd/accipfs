package core

// NodeManager ...
type NodeManager interface {
	Push(node Node)
	Store() error
	Load() error
}
