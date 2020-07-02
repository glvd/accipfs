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

//func parseIPFSAddrInfo(b []byte) (peer.AddrInfo, error) {
//	var info jsonIPFSAddrInfo
//	err := json.Unmarshal(b, &info)
//	if err != nil {
//		return peer.AddrInfo{}, err
//	}
//	var addrs []ma.Multiaddr
//	for i := range info.Addrs {
//		multiaddr, err := ma.NewMultiaddr(info.Addrs[i])
//		if err != nil {
//			continue
//		}
//		addrs = append(addrs, multiaddr)
//	}
//	fromString, err := peer.IDFromString(info.ID)
//	if err != nil {
//		return peer.AddrInfo{}, err
//	}
//	return peer.AddrInfo{
//		ID:    fromString,
//		Addrs: addrs,
//	}, nil
//}

func parseAddrInfo(b []byte, addrInfo *AddrInfo) error {
	var info jsonAddrInfo
	err := json.Unmarshal(b, &info)
	if err != nil {
		return fmt.Errorf("unmarshal address info failed:%w", err)

	}
	addrInfo.ID = info.ID
	//if info.IPFSAddrInfo != "" {
	//	err = json.Unmarshal([]byte(info.IPFSAddrInfo), &addrInfo.IPFSAddrInfo)
	//	if err != nil {
	//		return fmt.Errorf("unmarshal ipfs address info failed:%w", err)
	//	}
	//}
	//ipfsAddrInfo, err := parseIPFSAddrInfo([]byte(info.IPFSAddrInfo))
	//if err != nil {
	//	return err
	//}
	//addrInfo.IPFSAddrInfo = ipfsAddrInfo
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
func (info *AddrInfo) MarshalJSON() ([]byte, error) {
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
