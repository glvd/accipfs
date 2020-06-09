package core

// NodeManager ...
type NodeManager interface {
	Push(node Node)
	Load() error
}
