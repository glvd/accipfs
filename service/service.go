package service

import (
	"github.com/glvd/accipfs/config"
	"sync"
)

// Service ...
type Service struct {
	once *sync.Once
	cfg  *config.Config
	i    *nodeIPFS
	e    *nodeETH
}

// New ...
func New(config config.Config) (s *Service, e error) {
	s = &Service{
		cfg:  &config,
		once: &sync.Once{},
	}
	s.i, e = newNodeIPFS(config)
	if e != nil {
		return nil, e
	}
	s.e, e = newETH(config)
	if e != nil {
		return nil, e
	}
	return s, e
}

// Run ...
func (s *Service) Run() {
	s.once.Do(func() {
	})
}
