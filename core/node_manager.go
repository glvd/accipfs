package core

import (
	"encoding/json"
	"net"
	"sync"
)

// LocalData ...
type LocalData struct {
	lock sync.RWMutex
	Node NodeInfo
}

// Marshal ...
func (l *LocalData) Marshal() ([]byte, error) {
	l.lock.RLock()
	defer l.lock.Unlock()
	marshal, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}
	return marshal, err
}

// Unmarshal ...
func (l *LocalData) Unmarshal(bytes []byte) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	return json.Unmarshal(bytes, l)
}

// JSON ...
func (l *LocalData) JSON() string {
	marshal, err := json.Marshal(l)
	if err != nil {
		return ""
	}
	return string(marshal)
}

// Sync ...
func (l *LocalData) Sync(data *LocalData) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.Node = data.Node
}

// NodeManager ...
type NodeManager interface {
	NodeAPI
	Local() LocalData
	Close()
	Push(n Node)
	Range(f func(key string, node Node) bool)
	Conn(c net.Conn) (Node, error)
	Store() error
	Load() error
}
