package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"net"
	"net/http"
)

// API ...
type API struct {
	cfg        *config.Config
	eng        *gin.Engine
	listener   net.Listener
	serv       *http.Server
	controller *Controller
}

// Ping ...
func (a *API) Ping(req *core.PingReq) (*core.PingResp, error) {
	panic("implement me")
}

// ID ...
func (a *API) ID(req *core.IDReq) (*core.IDResp, error) {
	panic("implement me")
}

// New ...
func newAPI(cfg *config.Config, controller *Controller) core.API {
	return &API{
		cfg:        cfg,
		controller: controller,
		eng:        gin.Default(),
		serv:       &http.Server{},
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
	if a.cfg.API.UseTLS {
		go a.serv.ServeTLS(l, a.cfg.API.TLS.KeyFile, a.cfg.API.TLS.KeyPassFile)
		return nil
	}
	go a.serv.Serve(l)
	return nil
}

// Stop ...
func (a *API) Stop() error {
	if a.serv != nil {
		return a.serv.Close()
	}
	return nil
}
