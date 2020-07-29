package controller

import (
	"context"
	"github.com/glvd/accipfs/config"
	files "github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"os"

	//"github.com/glvd/accipfs/contract/dtag"
	"github.com/glvd/accipfs/core"
	"github.com/libp2p/go-libp2p-core/peer"
	"go.uber.org/atomic"
	"sync"
	"time"
)

// ServiceIndex ...
type ServiceIndex int

const (
	// IndexETH ...
	IndexETH ServiceIndex = iota
	// IndexIPFS ...
	IndexIPFS

	// IndexMax ...
	IndexMax
)

// Controller ...
type Controller struct {
	isRunning *atomic.Bool
	services  []core.ControllerService
	ethNode   *nodeBinETH
	ipfsNode  *nodeLibIPFS
	cfg       *config.Config
}

// New ...
func New(cfg *config.Config) *Controller {
	c := &Controller{
		cfg:      cfg,
		services: make([]core.ControllerService, IndexMax),
	}

	if cfg.ETH.Enable {
		eth := newNodeBinETH(cfg)
		eth.MessageHandle(func(s string) {
			output("[info]", s)
			//log.Infow(s, "tag", "eth")
		})
		c.services[IndexETH] = eth
		c.ethNode = eth
	}
	if cfg.IPFS.Enable {
		ipfs := newNodeLibIPFS(cfg)
		ipfs.MessageHandle(func(s string) {
			output("[datastore]", s)
		})
		c.services[IndexIPFS] = ipfs
		c.ipfsNode = ipfs
	}
	c.isRunning = atomic.NewBool(false)
	return c
}

// Initialize ...
func (c *Controller) Initialize() (e error) {
	for _, service := range c.services {
		if service == nil {
			continue
		}
		e = service.Initialize()
		if e != nil {
			return e
		}
	}
	return
}

// WaitAllReady ...
func (c Controller) WaitAllReady() {
	for {
	Reset:
		time.Sleep(3 * time.Second)
		for _, service := range c.services {
			if service == nil {
				continue
			}
			if b := service.IsReady(); !b {
				goto Reset
			}
		}
		break
	}
}

// Run ...
func (c *Controller) Run() {
	if !c.isRunning.CAS(false, true) {
		return
	}
	wg := &sync.WaitGroup{}
	for idx := range c.services {
		if c.services[idx] != nil {
			wg.Add(1)
			go func(service core.ControllerService) {
				defer wg.Done()
				if err := service.Start(); err != nil {
					log.Errorw("controller start failed", "err", err)
				}
			}(c.services[idx])
		}
	}
	wg.Wait()
}

// Stop ...
func (c *Controller) Stop() (e error) {
	for i, service := range c.services {
		if err := service.Stop(); err != nil {
			//stop all and collect exceptions
			logE("stop error", "index", i, "error", err)
			e = err
		}
	}
	c.isRunning.Store(false)
	return
}

func (c *Controller) dataNode() *nodeLibIPFS {
	return c.ipfsNode
}

func (c *Controller) infoNode() *nodeBinETH {
	return c.ethNode
}

// ID ...
func (c *Controller) ID(ctx context.Context) (*core.DataStoreInfo, error) {
	return &core.DataStoreInfo{}, nil
}

// DataStoreAPI ...
func (c *Controller) DataStoreAPI() core.DataStoreAPI {
	return c
}

// Get ...
func (c *Controller) GetUnixfs(ctx context.Context, urlPath string, endpoint string) (node files.Node, err error) {
	parsedPath := path.New(urlPath)
	if err := parsedPath.IsValid(); err != nil {
		return nil, err
	}

	resolvedPath, err := c.dataNode().ResolvePath(ctx, parsedPath)
	if err != nil {
		return nil, err
	}
	node, err = c.dataNode().Unixfs().Get(ctx, resolvedPath)
	if err != nil {
		return nil, err
	}
	_, ok := node.(files.Directory)
	if ok {
		if endpoint != "" {
			node, err = c.dataNode().Unixfs().Get(context.TODO(), path.Join(resolvedPath, endpoint))
			if err != nil {
				return nil, err
			}
		}
	}
	return node, err
}

// PinLs ...
func (c *Controller) PinLs(req *core.DataStoreReq) (*core.DataStoreResp, error) {
	ls, err := c.dataNode().Pin().Ls(context.TODO(), func(settings *options.PinLsSettings) error {
		settings.Type = "recursive"
		return nil
	})
	if err != nil {
		return nil, err
	}

	var pins []string
	for v := range ls {
		if v.Err() != nil {
			return nil, v.Err()
		}
		log.Infow("show pins", "data", v.Path().Cid())
		pins = append(pins, v.Path().Cid().String())
	}
	return &core.DataStoreResp{Pins: pins}, nil
}

// HandleSwarm ...
func (c *Controller) HandleSwarm(info peer.AddrInfo) error {
	return c.ipfsNode.Swarm().Connect(context.TODO(), info)
}

// UploadFile ...
func (c *Controller) UploadFile(req *core.UploadReq) (*core.UploadResp, error) {
	stat, e := os.Stat(req.Path)
	if e != nil {
		return &core.UploadResp{}, e
	}
	var node files.Node
	//var err error
	if !stat.IsDir() {
		file, e := os.Open(req.Path)
		if e != nil {
			return &core.UploadResp{}, e
		}
		node = files.NewReaderFile(file)
	} else {
		sf, e := files.NewSerialFile(req.Path, false, stat)
		if e != nil {
			return &core.UploadResp{}, e
		}
		node = sf
	}
	opts := func(settings *options.UnixfsAddSettings) error {
		settings.Pin = true
		return nil
	}

	resolved, e := c.dataNode().Unixfs().Add(context.TODO(), node, opts)
	if e != nil {
		return &core.UploadResp{}, e
	}
	return &core.UploadResp{
		Hash: resolved.Cid().String(),
	}, nil
}
