package service

import (
	"github.com/glvd/accipfs/config"
	"net/http"
)

type httpHandle struct {
	cfg *config.Config
}

func newHTTPHandle(cfg *config.Config) (*httpHandle, error) {
	return &httpHandle{
		cfg: cfg,
	}, nil
}

// Handler ...
func (s *httpHandle) Handler() (string, http.Handler) {
	return "/api", nil
}
