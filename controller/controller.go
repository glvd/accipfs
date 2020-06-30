package controller

import (
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"sync"
)

// ServiceIndex ...
type ServiceIndex int

const (
	// IndexETH ...
	IndexETH ServiceIndex = iota
	// IndexIPFS ...
	IndexIPFS
	// IndexAPI ...
	IndexAPI
	// IndexMax ...
	IndexMax
)

// Controller ...
type Controller struct {
	wg       *sync.WaitGroup
	services []core.ControllerService
	api      core.API
}

// New ...
func New(cfg *config.Config) *Controller {
	c := &Controller{
		services: make([]core.ControllerService, IndexMax),
	}

	api := newAPI(cfg)
	if cfg.ETH.Enable {
		eth := newNodeBinETH(cfg)
		eth.MessageHandle(func(s string) {
			log.Infow(s, "tag", "eth")
		})
		c.services[IndexETH] = eth
		api.ethNode = eth
	}
	if cfg.IPFS.Enable {
		ipfs := newNodeBinIPFS(cfg)
		ipfs.MessageHandle(func(s string) {
			log.Infow(s, "tag", "ipfs")
		})
		c.services[IndexIPFS] = ipfs
		api.ipfsNode = ipfs
	}
	c.services[IndexAPI] = api
	c.api = api
	c.wg = &sync.WaitGroup{}
	return c
}

// Initialize ...
func (c *Controller) Initialize() (e error) {
	for _, service := range c.services {
		e = service.Initialize()
		if e != nil {
			return e
		}
	}
	return
}

// Run ...
func (c *Controller) Run() {
	wg := &sync.WaitGroup{}
	for idx := range c.services {
		wg.Add(1)
		go func(service core.ControllerService) {
			defer c.wg.Done()
			if err := service.Start(); err != nil {
				return
			}
		}(c.services[idx])
	}
	wg.Wait()
}

// StopRun ...
func (c *Controller) StopRun() (e error) {
	for i, service := range c.services {
		if err := service.Stop(); err != nil {
			//stop all and collect exceptions
			logE("stop error", "index", i, "error", err)
			e = err
		}
	}
	return
}

// LocalAPI ...
func (c *Controller) LocalAPI() core.API {
	return c.api
}
