package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/glvd/accipfs/config"
	"net/http"
	"strings"
)

type httpHandle struct {
	cfg    *config.Config
	eng    *gin.Engine
	linker *BustLinker
}

// Success ...
var Success = "success"

func newHTTPHandle(cfg *config.Config, linker *BustLinker, eng interface{}) (*httpHandle, error) {
	g, b := eng.(*gin.Engine)
	if !b {
		g = gin.New()
	}
	//g.Use(cors.New(cors.Config{
	//	AllowOrigins:     []string{"*"},
	//	AllowMethods:     []string{"POST", "GET", "PUT", "PATCH"},
	//	AllowHeaders:     []string{"Origin"},
	//	ExposeHeaders:    []string{"Content-Length"},
	//	AllowCredentials: true,
	//	AllowOriginFunc: func(origin string) bool {
	//		return origin == "*"
	//	},
	//	MaxAge: 12 * time.Hour,
	//}))
	g.Use(gin.Recovery(), func(context *gin.Context) {
		logI("output url", "url", context.Request.URL.String())
	})

	h := &httpHandle{
		cfg:    cfg,
		eng:    g,
		linker: linker,
	}
	h.handleList()
	return h, nil
}
func (s *httpHandle) handleList() {
	g := s.eng.Group("/api")
	g.GET("/ping", s.Ping())
	if s.cfg.Debug {
		g.GET("/debug", s.Debug())
	}

	v0 := g.Group("v0")
	v0.POST("/info", s.Info())
	v0.GET("/get", s.Get())
}

// Ping ...
func (s *httpHandle) Ping() func(context *gin.Context) {
	return func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"status":  Success,
		})
	}
}

// Info ...
func (s *httpHandle) Info() func(context *gin.Context) {
	return func(context *gin.Context) {
		var err error
		j := struct {
			No string
		}{}
		err = context.BindJSON(&j)
		if err != nil {
			failedResult(context, fmt.Errorf("bind json failed(%w)", err))
			return
		}
		var rs string
		err = s.linker.tagInfo(strings.ToUpper(j.No), &rs)
		if err != nil {
			failedResult(context, fmt.Errorf("address(%s):%w", s.cfg.ETH.DTagAddr, err))
			return
		}
		context.JSON(http.StatusOK, gin.H{
			"message": rs,
			"status":  Success,
		})
	}
}

// Get ...
func (s *httpHandle) Get() func(context *gin.Context) {
	return func(context *gin.Context) {
		context.Redirect(http.StatusMovedPermanently, spliceGetURL("api/v0/get"))
	}
}

// Debug ...
func (s *httpHandle) Debug() func(context *gin.Context) {
	return func(context *gin.Context) {
		uri := context.Query("uri")
		context.Redirect(http.StatusMovedPermanently, spliceGetURL(uri))
	}
}

func spliceGetURL(uri string) string {
	return fmt.Sprintf("%s/%s", config.IPFSAddrHTTP(), uri)
}

// Handler ...
func (s *httpHandle) Handler() (string, http.Handler) {
	return "/api/*uri", s
}

// ServeHTTP ...
func (s *httpHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.eng.ServeHTTP(w, r)
}

func failedResult(ctx *gin.Context, err error) (b bool) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "failed",
		"error":  err.Error(),
	})
	return
}
