package core

import "time"

// Node ...
type Node struct {
	Addr
	NodeInfo
	LastTime        time.Time
	ProtocolVersion string
}
