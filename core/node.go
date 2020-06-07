package core

import "time"

// Node ...
type Node struct {
	ID   string
	Addr []Addr
	NodeInfo
	LastTime        time.Time
	ProtocolVersion string
}
