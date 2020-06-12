package api

import (
	"github.com/gin-gonic/gin"
	"github.com/glvd/accipfs/config"
	"net"
	"net/http"
)

// API ...
type API struct {
	cfg      *config.Config
	eng      *gin.Engine
	listener net.Listener
	serv     *http.Server
}

// New ...
func New(cfg *config.Config) *API {
	return &API{
		cfg:  cfg,
		eng:  gin.Default(),
		serv: &http.Server{},
	}
}

// Start ...
func (a *API) Start() error {
	l, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.IPv4zero,
		Port: a.cfg.API.Port,
	})
	if err != nil {
		return err
	}
	go a.serv.ServeTLS(l, a.cfg.API.TLS.KeyFile, a.cfg.API.TLS.KeyPassFile)
	return nil
}

// Stop ...
func (a *API) Stop() error {
	if a.serv != nil {
		return a.serv.Close()
	}
	return nil
}
