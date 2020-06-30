package core

import (
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

// AddrInfo ...
type AddrInfo struct {
	ID           string
	Addrs        map[ma.Multiaddr]bool
	IPFSAddrInfo peer.AddrInfo
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

// SetIPFSAddrInfo ...
func (info *AddrInfo) SetIPFSAddrInfo(ipfsAddrInfo peer.AddrInfo) {
	info.IPFSAddrInfo = ipfsAddrInfo
}

// SetID ...
func (info *AddrInfo) SetID(id string) {
	info.ID = id
}

// AppendAddr ...
func (info *AddrInfo) AppendAddr(multiaddr ma.Multiaddr) {
	info.Addrs[multiaddr] = true
}

// GetAddrs ...
func (info *AddrInfo) GetAddrs() (addrs []ma.Multiaddr) {
	for addr := range info.Addrs {
		addrs = append(addrs, addr)
	}
	return
}
