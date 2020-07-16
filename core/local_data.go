package core

import (
	"encoding/json"
	"sync"
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
	LDs         map[string]uint8 //ipfs linked data
	Addrs       []string
	LastUpdate  int64
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
	l.lock.Lock()
	data = l.data
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
