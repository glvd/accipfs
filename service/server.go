package service

// NodeServer ...
type NodeServer interface {
	Start() error
	Init() error
}

// Server ...
type Server struct {
}
