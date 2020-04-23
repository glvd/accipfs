package service

import (
	"github.com/gin-gonic/gin"
	"github.com/glvd/accipfs/config"
	"net/http"
)

type httpHandle struct {
	cfg *config.Config
	eng *gin.Engine
}

func newHTTPHandle(cfg *config.Config) (*httpHandle, error) {
	g := gin.Default()

	h := &httpHandle{
		cfg: cfg,
		eng: g,
	}
	h.handleList()
	return h, nil
}

// Handler ...
func (s *httpHandle) Handler() (string, http.Handler) {
	return "/api", s.eng
}

func (s *httpHandle) handleList() {
	s.eng.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"status":  "success",
		})
	})
}
