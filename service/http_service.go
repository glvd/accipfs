package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/glvd/accipfs/config"
	"net/http"
)

type httpService struct {
	cfg    *config.Config
	server *http.Server
	route  *gin.Engine
}

func newHTTPService(cfg *config.Config) *httpService {
	s := &httpService{
		cfg:   cfg,
		route: gin.Default(),
	}
	port := fmt.Sprintf(":%d", s.cfg.Port)

	s.server = &http.Server{Addr: port, Handler: s.route}
	s.route.Any("/", func(c *gin.Context) {
		c.String(http.StatusOK, "text/plain; charset=utf-8", "service is already running")
	})
	return s
}

// RegisterHandle ...
func (s *httpService) Register(path string, handler http.Handler) error {
	s.route.Any(path, func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	})
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

// ServeHTTP ...
func (s *httpService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//TODO handle api
}
