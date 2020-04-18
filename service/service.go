package service

import (
	"time"

	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/controller"
)

// Service ...
type service struct {
	linker     *BustLinker
	server     *rpcServer
	controller *controller.Controller
}

// Service ...
type Service interface {
	Start() (e error)
	Stop() (e error)
}

// New ...
func New(cfg *config.Config) (s Service, e error) {
	linker, e := NewBustLinker(cfg)
	if e != nil {
		return nil, e
	}

	server, e := newRPCServer(cfg, linker)
	if e != nil {
		return nil, e
	}
	s = &service{
		controller: controller.New(cfg),
		linker:     linker,
		server:     server,
	}
	return s, nil
}

// Start ...
func (s *service) Start() error {
	s.controller.Run()

	go s.linker.Start()

	s.linker.WaitingForReady()

	var idError error
	for i := 0; i < 5; i++ {
		id, err := s.linker.localID()
		idError = err
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}
		s.linker.id = id
		break
	}

	if idError != nil {
		return idError
	}

	s.server.Start()
	return nil
}

// Stop ...
func (s *service) Stop() error {
	if err := s.server.Stop(); err != nil {
		return err
	}

	s.linker.Stop()

	return s.controller.StopRun()
}
