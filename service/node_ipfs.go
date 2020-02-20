package service

import (
	"context"
	"fmt"
	"github.com/glvd/accipfs"
	"github.com/glvd/accipfs/config"
	"github.com/ipfs/go-ipfs-http-client"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"os/exec"
	"path/filepath"
)

const ipfsPath = ".ipfs"
const ipfsAPI = "api"

type nodeClientIPFS struct {
	cfg config.Config
	api *httpapi.HttpApi
}

type nodeServerIPFS struct {
	name string
	cmd  *exec.Cmd
}

// PeerID ...
type PeerID struct {
	Addresses       []string `json:"Addresses"`
	AgentVersion    string   `json:"AgentVersion"`
	ID              string   `json:"ID"`
	ProtocolVersion string   `json:"ProtocolVersion"`
	PublicKey       string   `json:"PublicKey"`
}

func newNodeIPFS(config config.Config) (*nodeClientIPFS, error) {
	api, e := httpapi.NewPathApi(filepath.Join(config.Path, ipfsPath, ipfsAPI))
	if e != nil {
		return nil, e
	}
	return &nodeClientIPFS{
		cfg: config,
		api: api,
	}, nil
}

// SwarmConnect ...
func (i *nodeClientIPFS) SwarmConnect(ctx context.Context, addr string) (e error) {
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

func (i *nodeClientIPFS) connect() (e error) {
	ma, err := multiaddr.NewMultiaddr(i.cfg.IPFS.Addr)
	if err != nil {
		return err
	}
	i.api, e = httpapi.NewApi(ma)
	return
}

// SwarmPeers ...
func (i *nodeClientIPFS) SwarmPeers(ctx context.Context) ([]iface.ConnectionInfo, error) {
	return i.api.Swarm().Peers(ctx)
}

// ID get self node info
func (i *nodeClientIPFS) ID(ctx context.Context) (pid *PeerID, e error) {
	pid = &PeerID{}
	e = i.api.Request("id").Exec(ctx, pid)
	if e != nil {
		return nil, e
	}
	return pid, nil
}

// PinAdd ...
func (i *nodeClientIPFS) PinAdd(ctx context.Context, hash string) (e error) {
	p := path.New(hash)
	return i.api.Pin().Add(ctx, p)
}

// PinLS ...
func (i *nodeClientIPFS) PinLS(ctx context.Context) (pins []iface.Pin, e error) {
	return i.api.Pin().Ls(ctx)
}

// PinRm ...
func (i *nodeClientIPFS) PinRm(ctx context.Context, hash string) (e error) {
	p := path.New(hash)
	return i.api.Pin().Rm(ctx, p)
}

// Start ...
func (n *nodeServerIPFS) Start() {
	fmt.Println("starting", n.name)
}

// NodeServerIPFS ...
func NodeServerIPFS(cfg config.Config) Node {
	cmd := exec.Command(cfg.IPFS.Name, "")
	cmd.Env = accipfs.Environ()
	return &nodeServerIPFS{cmd: cmd}
}
