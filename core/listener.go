package core

// Listener ...
type Listener interface {
	Listen() error
	Stop() error
}
