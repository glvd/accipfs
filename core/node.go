package core

import "time"

// Node ...
type Node struct {
	NodeAddress
	NodeInfo
	LastTime        time.Time
	ProtocolVersion string
}
