package service

import (
	"github.com/glvd/accipfs/config"
	"net/http"
)

// Empty ...
type Empty struct {
}

// Accelerate ...
type Accelerate struct {
	cfg config.Config
}

// Ping ...
func (n *Accelerate) Ping(r *http.Request, s *Empty, result *string) error {
	*result = "pong pong pong"
	return nil
}

// Account ...
func (n *Accelerate) Account(r *http.Request, s *Empty, result *string) error {
	*result = n.cfg.Account
	return nil
}
