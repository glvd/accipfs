package core

import ma "github.com/multiformats/go-multiaddr"

// AddrInfo ...
type AddrInfo struct {
	ID    string
	Addrs map[ma.Multiaddr]bool
}

// Append ...
func (info *AddrInfo) Append(multiaddr ma.Multiaddr) {
	info.Addrs[multiaddr] = true
}
