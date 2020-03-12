package service

import (
	"github.com/glvd/accipfs/config"
)

// NodeServer ...
type NodeServer interface {
	Start() error
	Init() error
	Stop() error
}

// Server ...
type Server struct {
	cfg *config.Config
}

// NewServer ...
func NewServer(cfg config.Config) *Server {
	return &Server{cfg: &cfg}
}
