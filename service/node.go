package service

import "time"

// HandleInfo ...
type HandleInfo struct {
	ServiceName string
	Data        interface{}
	Callback    HandleCallback
}

// HandleCallback ...
type HandleCallback func(src interface{})

// Node ...
type Node interface {
	Start()
}

var dateKey = time.Date(2019, time.November, 11, 10, 20, 10, 300, time.Local)
