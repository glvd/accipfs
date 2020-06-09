package core

// ControllerService ...
type ControllerService interface {
	Start() error
	Stop() error
	Initialize() error
	MessageHandle(func(s string))
}
