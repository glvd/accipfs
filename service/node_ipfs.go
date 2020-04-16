package service

import (
	"context"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/ipfs/go-ipfs-http-client"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
)

const ipfsPath = ".ipfs"
const ipfsAPI = "api"

type nodeIPFS struct {
	*serviceNode
	cfg *config.Config
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

func newNodeIPFS(cfg *config.Config) (*nodeIPFS, error) {
	node := &nodeIPFS{
		cfg:         cfg,
		serviceNode: nodeInstance(),
	}
	if err := node.connect(); err != nil {
		return nil, err
	}
	return node, nil
}

// SwarmConnect ...
func (n *nodeIPFS) SwarmConnect(ctx context.Context, addr string) (e error) {
	ma, e := multiaddr.NewMultiaddr(addr)
	if e != nil {
		return e
	}
	info, e := peer.AddrInfoFromP2pAddr(ma)
	if e != nil {
		return e
	}
	e = n.api.Swarm().Connect(ctx, *info)
	if e != nil {
		return e
	}
	return nil
}

func (n *nodeIPFS) connect() (e error) {
	ma, err := multiaddr.NewMultiaddr(config.IPFSAddr())
	if err != nil {
		return err
	}
	n.api, e = httpapi.NewApi(ma)
	return
}

// SwarmPeers ...
func (n *nodeIPFS) SwarmPeers(ctx context.Context) ([]iface.ConnectionInfo, error) {
	return n.api.Swarm().Peers(ctx)
}

// ID get self serviceNode info
func (n *nodeIPFS) ID(ctx context.Context) (pid *core.DataStoreNode, e error) {
	pid = &core.DataStoreNode{}
	e = n.api.Request("id").Exec(ctx, pid)
	if e != nil {
		return nil, e
	}
	return pid, nil
}

// PinAdd ...
func (n *nodeIPFS) PinAdd(ctx context.Context, hash string) (e error) {
	p := path.New(hash)
	return n.api.Pin().Add(ctx, p, options.Pin.Recursive(true))
}

// PinLS ...
func (n *nodeIPFS) PinLS(ctx context.Context) (pins []iface.Pin, e error) {
	return n.api.Pin().Ls(ctx, options.Pin.Type.Recursive())
}

// PinRm ...
func (n *nodeIPFS) PinRm(ctx context.Context, hash string) (e error) {
	p := path.New(hash)
	return n.api.Pin().Rm(ctx, p)
}

// IsReady ...
func (n *nodeIPFS) IsReady() bool {
	ma, err := multiaddr.NewMultiaddr(config.IPFSAddr())
	if err != nil {
		return false
	}
	api, e := httpapi.NewApi(ma)
	if e != nil {
		logE("new serviceNode ipfs", "error", e)
		return false
	}
	n.api = api
	return true
}
