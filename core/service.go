package core

// ControllerService ...
type ControllerService interface {
	Start() error
	Stop() error
	Initialize() error
	IsReady() bool
	MessageHandle(func(s string))
}
