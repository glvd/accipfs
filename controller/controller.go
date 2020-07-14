package controller

import (
	"context"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/contract/dtag"
	"github.com/glvd/accipfs/core"
	"go.uber.org/atomic"
	"sync"
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

// New ...
func New(cfg *config.Config) *Controller {
	c := &Controller{
		cfg:      cfg,
		services: make([]core.ControllerService, IndexMax),
	}

	if cfg.ETH.Enable {
		eth := newNodeBinETH(cfg)
		eth.MessageHandle(func(s string) {
			log.Infow(s, "tag", "eth")
		})
		c.services[IndexETH] = eth
		c.ethNode = eth
	}
	if cfg.IPFS.Enable {
		ipfs := newNodeBinIPFS(cfg)
		ipfs.MessageHandle(func(s string) {
			log.Infow(s, "tag", "ipfs")
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

// PinLs ...
func (c *Controller) PinLs() ([]string, error) {
	ls, err := c.dataNode().PinLS(context.TODO())
	if err != nil {
		return nil, err
	}

	var pins []string
	for v := range ls {
		pins = append(pins, v.Path().String())
	}
	return pins, nil
}

//// API ...
//func (c *Controller) API() core.API {
//	return service.NewAPIContext(c.cfg, m, c)
//}
