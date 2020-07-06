package controller

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/libp2p/go-libp2p-core/peer"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"go.uber.org/atomic"
)

// Context ...
type Context struct {
	cfg        *config.Config
	eng        *gin.Engine
	listener   net.Listener
	serv       *http.Server
	ready      *atomic.Bool
	controller *Controller
	ethNode    *nodeBinETH
	ipfsNode   *nodeBinIPFS
	msg        func(s string)
}

// AddrInfo ...
func (c *Context) AddrInfo(req *core.AddrReq) (*core.AddrResp, error) {
	panic("todo:AddrInfo")
}

// Ping ...
func (c *Context) Ping(req *core.PingReq) (*core.PingResp, error) {
	return &core.PingResp{
		Data: "pong",
	}, nil
}

// ID ...
func (c *Context) ID(req *core.IDReq) (*core.IDResp, error) {
	//loadAccount, err := account.LoadAccount(c.cfg)
	//if err != nil {
	//	return nil, err
	//}
	//loadAccount.Identity
	//log.Infow("get id", "account", loadAccount)
	//fromString := peer.ID(c.cfg.Identity)
	//if err != nil {
	//	log.Errorw("id from string", "id", c.cfg.Identity, "err", err)
	//	return nil, err
	//}
	fromString, err := peer.Decode(c.cfg.Identity)
	if err != nil {
		return nil, err
	}
	log.Infow("get id", "id", fromString.String())
	key, err := fromString.ExtractPublicKey()
	if err != nil {
		return nil, err
	}
	bytes, err := key.Bytes()
	if err != nil {
		return nil, err
	}
	pubKey := base64.StdEncoding.EncodeToString(bytes)
	log.Infow("result id", "id", c.cfg.Identity, "key", pubKey)
	return &core.IDResp{
		Name:      c.cfg.Identity,
		PublicKey: pubKey,
	}, nil
}

// New ...
func newAPI(cfg *config.Config) *Context {
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	eng := gin.Default()
	return &Context{
		cfg:   cfg,
		eng:   eng,
		ready: atomic.NewBool(false),
		serv: &http.Server{
			Handler: eng,
		},
	}
}

// Start ...
func (c *Context) Start() error {
	l, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.IPv4zero,
		Port: c.cfg.API.Port,
	})
	if err != nil {
		return err
	}
	c.registerRoutes()
	if c.cfg.API.UseTLS {
		go c.serv.ServeTLS(l, c.cfg.API.TLS.KeyFile, c.cfg.API.TLS.KeyPassFile)
		return nil
	}
	go c.serv.Serve(l)
	c.ready.Store(true)
	return nil
}

func (c *Context) registerRoutes() {
	api := c.eng.Group("/api")
	api.GET("/ping", c.ping)
	if c.cfg.Debug {
		api.GET("/debug", c.debug)
	}

	v0 := api.Group(c.cfg.API.Version)
	v0.POST("/id", c.id)
	v0.GET("/get", c.get)
	v0.GET("/query", c.query)
}

// Stop ...
func (c *Context) Stop() error {
	if c.serv != nil {
		if err := c.serv.Shutdown(context.TODO()); err != nil {
			return err
		}
	}
	return nil
}

// Initialize ...
func (c *Context) Initialize() error {
	//nothing
	return nil
}

// IsReady ...
func (c *Context) IsReady() bool {
	return c.ready.Load()
}

// MessageHandle ...
func (c *Context) MessageHandle(f func(s string)) {
	if f != nil {
		c.msg = f
	}
}

func (c *Context) id(ctx *gin.Context) {
	id, err := c.ID(&core.IDReq{})
	JSON(ctx, id, err)
}

func (c *Context) get(ctx *gin.Context) {
	ctx.Redirect(http.StatusMovedPermanently, ipfsGetURL("api/v0/get"))
}

func ipfsGetURL(uri string) string {
	return fmt.Sprintf("%s/%s", config.IPFSAddrHTTP(), uri)
}

func (c *Context) ping(ctx *gin.Context) {
	ping, err := c.Ping(&core.PingReq{})
	JSON(ctx, ping, err)
}

func (c *Context) debug(ctx *gin.Context) {
	uri := ctx.Query("uri")
	ctx.Redirect(http.StatusFound, ipfsGetURL(uri))
}

func (c *Context) query(ctx *gin.Context) {
	var err error
	j := struct {
		No string
	}{}
	err = ctx.BindJSON(&j)
	if err != nil {
		JSON(ctx, "", fmt.Errorf("query failed(%w)", err))
		return
	}
	dTag, e := c.ethNode.DTag()
	if e != nil {
		JSON(ctx, "", fmt.Errorf("query failed(%w)", e))
		return
	}
	message, e := dTag.GetTagMessage(&bind.CallOpts{Pending: true}, "video", j.No)
	if e != nil {
		JSON(ctx, "", fmt.Errorf("query failed(%w)", e))
		return
	}

	if message.Size.Int64() > 0 {
		JSON(ctx, message.Value[0], nil)
		return
	}
	JSON(ctx, "", nil)
}

// JSON ...
func JSON(c *gin.Context, v interface{}, e error) {
	if e != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "failed",
			"error":  e.Error(),
		})
		return
	}
	m, e := json.Marshal(v)
	if e != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "failed",
			"error":  e.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": string(m),
	})

}
