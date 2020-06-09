package core

// ControllerService ...
type ControllerService interface {
	Start() error
	Stop() error
	Init() error
	MessageHandle(func(s string))
}
