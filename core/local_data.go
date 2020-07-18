package core

import (
	"encoding/json"
	"sync"
	"time"
)

// SafeLocalData ...
type SafeLocalData interface {
	JSONer
	JSON() string
	Update(f func(data *LocalData))
	Data() (data LocalData)
}

// SafeLocalData ...
type safeLocalData struct {
	lock sync.RWMutex
	data LocalData
}

// LocalData ...
type LocalData struct {
	Initialized bool
	Node        NodeInfo
	Nodes       map[string]NodeInfo
	LDs         map[string]uint8 //ipfs linked data
	Addrs       []string
	LastUpdate  int64
}

// DefaultLocalData ...
func DefaultLocalData() *LocalData {
	return &LocalData{
		LDs:        make(map[string]uint8),
		LastUpdate: time.Now().Unix(),
		Nodes:      make(map[string]NodeInfo),
	}
}

// Marshal ...
func (l *safeLocalData) Marshal() ([]byte, error) {
	l.lock.RLock()
	marshal, err := json.Marshal(l.data)
	l.lock.RUnlock()
	if err != nil {
		return nil, err
	}
	return marshal, err
}

// Unmarshal ...
func (l *safeLocalData) Unmarshal(bytes []byte) (err error) {
	l.lock.Lock()
	err = json.Unmarshal(bytes, &l.data)
	l.lock.Unlock()
	return
}

// JSON ...
func (l *safeLocalData) JSON() string {
	marshal, err := l.Marshal()
	if err != nil {
		return ""
	}
	return string(marshal)
}

// Update ...
func (l *safeLocalData) Update(f func(data *LocalData)) {
	l.lock.Lock()
	f(&l.data)
	l.lock.Unlock()
}

// Data ...
func (l *safeLocalData) Data() (data LocalData) {
	data = LocalData{
		Initialized: false,
		Node:        NodeInfo{},
		Nodes:       make(map[string]NodeInfo),
		LDs:         make(map[string]uint8),
		Addrs:       nil,
		LastUpdate:  0,
	}
	l.lock.Lock()
	if l.data.Addrs != nil {
		copy(data.Addrs, l.data.Addrs)
	}
	for s := range l.data.Nodes {
		data.Nodes[s] = l.data.Nodes[s]
	}
	for s := range l.data.LDs {
		data.LDs[s] = l.data.LDs[s]
	}
	l.lock.Unlock()
	return
}

// Safe ...
func (l LocalData) Safe() SafeLocalData {
	return &safeLocalData{
		lock: sync.RWMutex{},
		data: l,
	}
}
