package controller

import (
	"errors"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	ma "github.com/multiformats/go-multiaddr"
	mnet "github.com/multiformats/go-multiaddr-net"
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
	// IndexAPI ...
	IndexAPI
	// IndexMax ...
	IndexMax
)

// Controller ...
type Controller struct {
	manager   core.NodeManager
	isRunning *atomic.Bool
	services  []core.ControllerService
	api       core.API
}

// New ...
func New(cfg *config.Config, manager core.NodeManager) *Controller {
	c := &Controller{
		services: make([]core.ControllerService, IndexMax),
	}

	api := newAPI(cfg, func(tag core.RequestTag, v interface{}) error {
		m, b := v.(ma.Multiaddr)
		if !b {
			return errors.New("wrong type convert")
		}
		dial, err := mnet.Dial(m)
		if err != nil {
			return err
		}
		manager.Conn(dial)
		return nil
	})
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
	c.isRunning = atomic.NewBool(false)
	c.services[IndexAPI] = api
	c.api = api
	c.manager = manager
	//c.wg = &sync.WaitGroup{}
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

// GetAPI ...
func (c *Controller) GetAPI() core.API {
	return c.api
}
