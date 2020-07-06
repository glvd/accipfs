package controller

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	ma "github.com/multiformats/go-multiaddr"
	"net"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gin-gonic/gin"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	ic "github.com/libp2p/go-libp2p-core/crypto"
	peer "github.com/libp2p/go-libp2p-core/peer"
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

var _ core.API = &Context{}

// NodeAddrInfo ...
func (c *Context) NodeAddrInfo(req *core.AddrReq) (*core.AddrResp, error) {
	if req.ID == "" {
		return &core.AddrResp{}, nil
	}
	panic("implement me")
}

// Ping ...
func (c *Context) Ping(req *core.PingReq) (*core.PingResp, error) {
	return &core.PingResp{
		Data: "pong",
	}, nil
}

// ID ...
func (c *Context) ID(req *core.IDReq) (*core.IDResp, error) {
	fromStringID, err := peer.Decode(c.cfg.Identity)
	if err != nil {
		return nil, err
	}
	log.Infow("get id", "id", fromStringID.String())
	pkb, err := base64.StdEncoding.DecodeString(c.cfg.PrivateKey)
	if err != nil {
		return nil, err
	}
	privateKey, err := ic.UnmarshalPrivateKey(pkb)
	if err != nil {
		return nil, err
	}
	publicKey := privateKey.GetPublic()
	bytes, err := publicKey.Bytes()
	if err != nil {
		return nil, err
	}
	pubString := base64.StdEncoding.EncodeToString(bytes)
	log.Infow("result id", "id", c.cfg.Identity, "public key", pubString)
	ipfsID, err := c.ipfsNode.ID(context.TODO())
	if err != nil {
		return nil, err
	}
	var multiAddress []ma.Multiaddr
	for _, address := range ipfsID.Addresses {
		multiaddr, err := ma.NewMultiaddr(address)
		if err != nil {
			continue
		}
		multiAddress = append(multiAddress, multiaddr)
	}
	return &core.IDResp{
		Name:      c.cfg.Identity,
		PublicKey: pubString,
		Addrs:     nil,
		DataStore: *ipfsID,
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
