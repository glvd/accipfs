package core

import ma "github.com/multiformats/go-multiaddr"

// AddrInfo ...
type AddrInfo struct {
	id    string
	addrs map[ma.Multiaddr]bool
}

// NewAddrInfo ...
func NewAddrInfo(id string, addrs ...ma.Multiaddr) *AddrInfo {
	_addrs := make(map[ma.Multiaddr]bool, len(addrs))
	for _, addr := range addrs {
		_addrs[addr] = true
	}
	return &AddrInfo{
		id:    id,
		addrs: _addrs,
	}
}

// SetID ...
func (info *AddrInfo) SetID(id string) {
	info.id = id
}

// ID ...
func (info *AddrInfo) ID() string {
	return info.id
}

// Append ...
func (info *AddrInfo) Append(multiaddr ma.Multiaddr) {
	info.addrs[multiaddr] = true
}

// Addrs ...
func (info *AddrInfo) Addrs() (addrs []ma.Multiaddr) {
	for addr := range info.addrs {
		addrs = append(addrs, addr)
	}
	return
}
