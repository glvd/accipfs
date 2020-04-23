package service

import (
	"context"
	"fmt"
	"github.com/glvd/accipfs/config"
	"github.com/gorilla/mux"
	"net/http"
)

type httpService struct {
	cfg    *config.Config
	server *http.Server
	route  *mux.Router
}

func newHTTPService(cfg *config.Config) *httpService {
	s := &httpService{
		cfg:   cfg,
		route: mux.NewRouter(),
	}
	port := fmt.Sprintf(":%d", s.cfg.Port)

	s.server = &http.Server{Addr: port, Handler: s.route}
	return s
}

// RegisterHandle ...
func (s *httpService) Register(path string, handler http.Handler) error {
	s.route.Handle(path, handler)
	return nil
}

// Start ...
func (s *httpService) Start() {
	output("JSON RPC service listen and serving on port", s.cfg.Port)
	if err := s.server.ListenAndServe(); err != nil {
		return
	}
}

// Stop ...
func (s *httpService) Stop() {
	if err := s.server.Shutdown(context.TODO()); err != nil {
		return
	}
}
