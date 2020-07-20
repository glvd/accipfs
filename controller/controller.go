package controller

import (
	"context"
	"fmt"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/contract/dtag"
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
	ipfsNode  *nodeBinIPFS
	cfg       *config.Config
}

// UploadFile ...
func (c *Controller) UploadFile(req *core.UploadReq) (*core.UploadResp, error) {
	fmt.Println("dummy:called on client")
	return &core.UploadResp{}, nil
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
		ipfs := newNodeBinIPFS(cfg)
		ipfs.MessageHandle(func(s string) {
			output("[datastore]", s)
			//log.Infow(s, "tag", "ipfs")
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
					return
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

func (c *Controller) dataNode() *nodeBinIPFS {
	return c.ipfsNode
}

func (c *Controller) infoNode() *nodeBinETH {
	return c.ethNode
}

// ID ...
func (c *Controller) ID(ctx context.Context) (*core.DataStoreInfo, error) {
	return c.dataNode().ID(ctx)
}

// DTag ...
func (c *Controller) DTag() (*dtag.DTag, error) {
	return c.infoNode().DTag()
}

// DataStoreAPI ...
func (c *Controller) DataStoreAPI() core.DataStoreAPI {
	return c
}

// PinLs ...
func (c *Controller) PinLs(req *core.DataStoreReq) (*core.DataStoreResp, error) {
	log.Infow("pin list request")
	ls, err := c.dataNode().PinLS(context.TODO())
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
	return c.ipfsNode.SwarmConnect(context.TODO(), info)
}
