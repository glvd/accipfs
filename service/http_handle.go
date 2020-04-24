package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/glvd/accipfs/config"
	"net/http"
)

type httpHandle struct {
	cfg *config.Config
	eng *gin.Engine
}

func newHTTPHandle(cfg *config.Config, eng interface{}) (*httpHandle, error) {
	g.Use(func(context *gin.Context) {
		logI("output url", "url", context.Request.URL.String())
	})

	g, b := eng.(*gin.Engine)
	if !b {
		return nil, fmt.Errorf("wrong gin type")
	}
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
	g := s.eng.Group("/api")
	g.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"status":  "success",
		})
	})
}
