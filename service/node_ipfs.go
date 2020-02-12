package service

import (
	"context"
	"github.com/glvd/accipfs/config"
	"github.com/ipfs/go-ipfs-http-client"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
)

type nodeIPFS struct {
	cfg config.Config
	api *httpapi.HttpApi
}

// PeerID ...
type PeerID struct {
	Addresses       []string `json:"Addresses"`
	AgentVersion    string   `json:"AgentVersion"`
	ID              string   `json:"ID"`
	ProtocolVersion string   `json:"ProtocolVersion"`
	PublicKey       string   `json:"PublicKey"`
}

func newNodeIPFS(config config.Config) (*nodeIPFS, error) {
	return &nodeIPFS{
		cfg: config,
	}, nil
}

func (i *nodeIPFS) SwarmConnect(ctx context.Context, addr string) (e error) {
	ma, e := multiaddr.NewMultiaddr(addr)
	if e != nil {
		return e
	}
	info, e := peer.AddrInfoFromP2pAddr(ma)
	if e != nil {
		return e
	}
	e = i.api.Swarm().Connect(ctx, *info)
	if e != nil {
		return e
	}
	return nil
}

func (i *nodeIPFS) connect() (e error) {
	ma, err := multiaddr.NewMultiaddr(i.cfg.IPFS.Addr)
	if err != nil {
		return err
	}
	i.api, e = httpapi.NewApi(ma)
	return
}

func (i *nodeIPFS) SwarmPeers(ctx context.Context) ([]iface.ConnectionInfo, error) {
	return i.api.Swarm().Peers(ctx)
}

// ID get self node info
func (i *nodeIPFS) ID(ctx context.Context) (pid *PeerID, e error) {
	pid = &PeerID{}
	e = i.api.Request("id").Exec(ctx, pid)
	if e != nil {
		return nil, e
	}
	return pid, nil
}
