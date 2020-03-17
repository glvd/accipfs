package service

import (
	"net/http"
)

// Empty ...
type Empty struct {
}

// Accelerate ...
type Accelerate struct {
}

// Ping ...
func (n *Accelerate) Ping(r *http.Request, s *Empty, result *string) error {
	*result = "pong pong pong"
	return nil
}
