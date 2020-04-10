package core

// Service ...
type Service interface {
	Start() error
	Stop() error
	Init() error
}
