package controller

import (
	"github.com/glvd/accipfs/config"
	"sync"
)

// ServiceIndex ...
type ServiceIndex int

const (
	// IndexETH ...
	IndexETH ServiceIndex = iota
	// IndexIPFS ...
	IndexIPFS
)

// Service ...
type Service interface {
	Start() error
	Stop() error
	Init() error
}

// Controller ...
type Controller struct {
	wg       *sync.WaitGroup
	services map[ServiceIndex]Service
}

// New ...
func New(cfg *config.Config) *Controller {
	c := &Controller{
		services: map[ServiceIndex]Service{},
	}
	c.wg = &sync.WaitGroup{}
	c.services[IndexETH] = newNodeBinETH(cfg)
	c.services[IndexIPFS] = newNodeBinIPFS(cfg)
	return c
}

// Init ...
func (c *Controller) Init() (e error) {
	for _, service := range c.services {
		e = service.Init()
		if e != nil {
			return e
		}
	}
	return
}

// Run ...
func (c *Controller) Run() {
	for _, service := range c.services {
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			if err := service.Start(); err != nil {
				return
			}
		}()
	}
	c.wg.Wait()
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
