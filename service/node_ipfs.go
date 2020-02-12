package service

import (
	"context"
	"github.com/glvd/accipfs/config"
	"github.com/ipfs/go-ipfs-http-client"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
)

type nodeIPFS struct {
	cfg config.Config
	api *httpapi.HttpApi
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
