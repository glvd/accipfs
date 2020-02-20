package service

import (
	"github.com/glvd/accipfs/config"
	"sync"
)

// Service ...
type Service struct {
	cfg        *config.Config
	serveMutex sync.RWMutex
	serve      []Node
	i          *nodeClientIPFS
	e          *nodeClientETH
}

// New ...
func New(config config.Config) (s *Service, e error) {
	s = &Service{
		cfg: &config,
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

// RegisterServer ...
func (s *Service) RegisterServer(node Node) {
	s.serveMutex.Lock()
	defer s.serveMutex.Unlock()
	s.serve = append(s.serve, node)
}

// Run ...
func (s *Service) Run() {
	s.serveMutex.RLock()
	defer s.serveMutex.RUnlock()
	for _, s := range s.serve {
		s.Start()
	}
}
