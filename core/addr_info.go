package core

import (
	"encoding/json"
	"fmt"
	ma "github.com/multiformats/go-multiaddr"
)

// AddrInfo ...
type AddrInfo struct {
	ID        string                `json:"id"`
	PublicKey string                `json:"public_key"`
	Addrs     map[ma.Multiaddr]bool `json:"addrs"`
	DataStore DataStoreInfo         `json:"data_store"`
	//IPFSAddrInfo peer.AddrInfo
}
type jsonIPFSAddrInfo struct {
	ID    string   `json:"id"`
	Addrs []string `json:"addrs"`
}
type jsonAddrInfo struct {
	ID        string        `json:"id"`
	PublicKey string        `json:"public_key"`
	Addrs     []string      `json:"addrs"`
	DataStore DataStoreInfo `json:"data_store"`
	//IPFSAddrInfo peer.AddrInfo `json:"ipfs_addr_info"`
}

func parseAddrInfo(b []byte, addrInfo *AddrInfo) error {
	var info jsonAddrInfo
	err := json.Unmarshal(b, &info)
	if err != nil {
		return fmt.Errorf("unmarshal address info failed:%w", err)

	}
	addrInfo.ID = info.ID
	addrInfo.PublicKey = info.PublicKey
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
		ID:        info.ID,
		PublicKey: info.PublicKey,
		Addrs:     nil,
	}
	for multiaddr := range info.Addrs {
		addrInfo.Addrs = append(addrInfo.Addrs, multiaddr.String())
	}
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

// SetDataStoreInfo ...
func (info *AddrInfo) SetDataStoreInfo(dataStoreInfo DataStoreInfo) {
	info.DataStore = dataStoreInfo
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
