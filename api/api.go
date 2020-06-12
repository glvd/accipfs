package api

import (
	"github.com/gin-gonic/gin"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"net"
	"net/http"
)

type api struct {
	route *gin.RouterGroup
}

// New ...
func New(cfg *config.Config) (core.API, error) {
	l, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.IPv4zero,
		Port: cfg.API.Port,
	})
	if err != nil {
		return nil, err
	}

}

// ServeHTTP ...
func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
