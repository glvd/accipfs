package core

import (
	"encoding/json"
	"sync"
)

// LocalDataLocker ...
type LocalDataLocker struct {
	lock sync.RWMutex
	data LocalData
}

// LocalData ...
type LocalData struct {
	Node NodeInfo
}

// Marshal ...
func (l *LocalDataLocker) Marshal() ([]byte, error) {
	l.lock.RLock()
	marshal, err := json.Marshal(l.data)
	l.lock.Unlock()
	if err != nil {
		return nil, err
	}
	return marshal, err
}

// Unmarshal ...
func (l *LocalDataLocker) Unmarshal(bytes []byte) (err error) {
	l.lock.Lock()
	err = json.Unmarshal(bytes, &l.data)
	l.lock.Unlock()
	return
}

// JSON ...
func (l *LocalDataLocker) JSON() string {
	marshal, err := l.Marshal()
	if err != nil {
		return ""
	}
	return string(marshal)
}

// Update ...
func (l *LocalDataLocker) Update(f func(data *LocalData)) {
	l.lock.Lock()
	f(&l.data)
	l.lock.Unlock()
}

// Data ...
func (l *LocalDataLocker) Data() (data LocalData) {
	l.lock.Lock()
	data = l.data
	l.lock.Unlock()
	return
}
