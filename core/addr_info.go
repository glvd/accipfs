package core

import (
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

// AddrInfo ...
type AddrInfo struct {
	ID           string
	PublicKey    string
	Addrs        map[ma.Multiaddr]bool
	IPFSAddrInfo peer.AddrInfo
}
type jsonIPFSAddrInfo struct {
	ID    string   `json:"id"`
	Addrs []string `json:"addrs"`
}
type jsonAddrInfo struct {
	ID           string        `json:"id"`
	Addrs        []string      `json:"addrs"`
	IPFSAddrInfo peer.AddrInfo `json:"ipfs_addr_info"`
}

func parseAddrInfo(b []byte, addrInfo *AddrInfo) error {
	var info jsonAddrInfo
	err := json.Unmarshal(b, &info)
	if err != nil {
		return fmt.Errorf("unmarshal address info failed:%w", err)

	}
	addrInfo.ID = info.ID
	addrs := make(map[ma.Multiaddr]bool, len(info.Addrs))
	for i := range info.Addrs {
		multiaddr, err := ma.NewMultiaddr(info.Addrs[i])
		if err != nil {
			continue
		}
		addrs[multiaddr] = true
	}
	addrInfo.Addrs = addrs
	return nil
}

// MarshalJSON ...
func (info AddrInfo) MarshalJSON() ([]byte, error) {
	addrInfo := jsonAddrInfo{
		ID:    info.ID,
		Addrs: nil,
	}
	for multiaddr := range info.Addrs {
		addrInfo.Addrs = append(addrInfo.Addrs, multiaddr.String())
	}
	//_, err := info.IPFSAddrInfo.MarshalJSON()
	//if err != nil {
	//	return nil, err
	//}
	//addrInfo.IPFSAddrInfo = string(v)
	return json.Marshal(addrInfo)
}

// UnmarshalJSON ...
func (info *AddrInfo) UnmarshalJSON(bytes []byte) error {
	return parseAddrInfo(bytes, info)
}

// NewAddrInfo ...
func NewAddrInfo(id string, addrs ...ma.Multiaddr) *AddrInfo {
	_addrs := make(map[ma.Multiaddr]bool, len(addrs))
	for _, addr := range addrs {
		_addrs[addr] = true
	}
	return &AddrInfo{
		ID:    id,
		Addrs: _addrs,
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
