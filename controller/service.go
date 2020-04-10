package controller

import "github.com/glvd/accipfs/config"

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
	services map[ServiceIndex]Service
}

// New ...
func New(cfg *config.Config) *Controller {
	c := &Controller{
		services: map[ServiceIndex]Service{},
	}
	c.services[IndexETH] = newNodeBinETH(cfg)
	c.services[IndexIPFS] = newNodeBinIPFS(cfg)
	return c
}

// Run ...
func (c *Controller) Run() {

}
