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

func newHTTPHandle(cfg *config.Config, eng interface{}) (*httpHandle, error) {
	g, b := eng.(*gin.Engine)
	if !b {
		g = gin.Default()
	}

	g.Use(func(context *gin.Context) {
		logI("output url", "url", context.Request.URL.String())
	})

	h := &httpHandle{
		cfg: cfg,
		eng: g,
	}
	h.handleList()
	return h, nil
}

// Handler ...
func (s *httpHandle) Handler() (string, http.Handler) {
	return "/api/*path", s.eng
}

func (s *httpHandle) handleList() {
	g := s.eng.Group("/api")
	g.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"status":  "success",
		})
	})
}
