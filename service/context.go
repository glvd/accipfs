package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/glvd/accipfs/controller"
	ma "github.com/multiformats/go-multiaddr"
	"net"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gin-gonic/gin"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	ic "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"go.uber.org/atomic"
)

// APIContext ...
type APIContext struct {
	cfg      *config.Config
	eng      *gin.Engine
	listener net.Listener
	serv     *http.Server
	ready    *atomic.Bool
	c        *controller.Controller
	m        core.NodeManager
	msg      func(s string)
}

var _ core.API = &APIContext{}

// NewAPIContext ...
func NewAPIContext(cfg *config.Config, m core.NodeManager, c *controller.Controller) *APIContext {
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	eng := gin.Default()
	return &APIContext{
		cfg:   cfg,
		eng:   eng,
		m:     m,
		c:     c,
		ready: atomic.NewBool(false),
		serv: &http.Server{
			Handler: eng,
		},
	}
}

// API ...
func (c *APIContext) API() core.API {
	return c
}

// NodeAddrInfo ...
func (c *APIContext) NodeAddrInfo(req *core.AddrReq) (*core.AddrResp, error) {
	id, err := c.ID(&core.IDReq{})
	if err != nil {
		return nil, err
	}
	info := core.NewAddrInfo(id.ID, id.Addrs...)
	info.PublicKey = id.PublicKey
	info.DataStore = id.DataStore
	//log.Infow("node addr", "info", info.DataStore, "data", id.DataStore)
	return &core.AddrResp{
		AddrInfo: *info,
	}, nil

}

// PinLs ...
func (c *APIContext) PinLs(req *core.DataStoreReq) (*core.DataStoreResp, error) {
	return c.DataStoreAPI().PinLs(req)
}

// DataStoreAPI ...
func (c *APIContext) DataStoreAPI() core.DataStoreAPI {
	return c.c
}

// Link ...
func (c *APIContext) Link(req *core.NodeLinkReq) (*core.NodeLinkResp, error) {
	return c.NodeAPI().Link(req)
}

// List ...
func (c *APIContext) List(req *core.NodeListReq) (*core.NodeListResp, error) {
	return c.NodeAPI().List(req)
}

// NodeAPI ...
func (c *APIContext) NodeAPI() core.NodeAPI {
	return c.m.NodeAPI()
}

// Add ...
func (c *APIContext) Add(req *core.AddReq) (*core.AddResp, error) {
	return c.NodeAPI().Add(req)
}

// Ping ...
func (c *APIContext) Ping(req *core.PingReq) (*core.PingResp, error) {
	return &core.PingResp{
		Data: "pong",
	}, nil
}

func (c *APIContext) nodeID() (string, string, error) {
	fromStringID, err := peer.Decode(c.cfg.Identity)
	if err != nil {
		return "", "", err
	}
	log.Infow("get id", "id", fromStringID.String())
	pkb, err := base64.StdEncoding.DecodeString(c.cfg.PrivateKey)
	if err != nil {
		return "", "", err
	}
	privateKey, err := ic.UnmarshalPrivateKey(pkb)
	if err != nil {
		return "", "", err
	}
	publicKey := privateKey.GetPublic()
	bytes, err := publicKey.Bytes()
	if err != nil {
		return "", "", err
	}
	pubString := base64.StdEncoding.EncodeToString(bytes)
	log.Infow("result id", "id", c.cfg.Identity, "public key", pubString)
	return fromStringID.Pretty(), pubString, nil
}

// ID ...
func (c *APIContext) ID(req *core.IDReq) (*core.IDResp, error) {
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

	ipfsID, err := c.c.ID(context.TODO())
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
	//log.Infow("result id", "id", c.cfg.Identity, "public key", pubString, "datastore", ipfsID)
	return &core.IDResp{
		ID:        c.cfg.Identity,
		PublicKey: pubString,
		Addrs:     nil,
		DataStore: *ipfsID,
	}, nil
}

// Start ...
func (c *APIContext) Start() error {
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

func (c *APIContext) registerRoutes() {
	c.eng.GET("/ping", c.ping)
	api := c.eng.Group("/api")
	if c.cfg.Debug {
		api.GET("/debug", c.debug)
	}

	v0 := api.Group(c.cfg.API.Version)
	v0.POST("/id", c.id)
	v0.POST("add", c.add())
	v0.POST("/node/link", c.nodeLink())
	v0.POST("/node/unlink", c.nodeUnlink())
	v0.POST("/node/list", c.nodeList())
	v0.POST("/datastore/pin/ls", c.datastorePinLs())
	v0.GET("/get", c.get)
	v0.GET("/query", c.query)
}

// Stop ...
func (c *APIContext) Stop() error {
	if c.serv != nil {
		if err := c.serv.Shutdown(context.TODO()); err != nil {
			return err
		}
	}
	return nil
}

// Initialize ...
func (c *APIContext) Initialize() error {
	//nothing
	return nil
}

// IsReady ...
func (c *APIContext) IsReady() bool {
	return c.ready.Load()
}

// MessageHandle ...
func (c *APIContext) MessageHandle(f func(s string)) {
	if f != nil {
		c.msg = f
	}
}

func (c *APIContext) setController(controller *controller.Controller) {
	c.c = controller
}

func (c *APIContext) id(ctx *gin.Context) {
	id, err := c.ID(&core.IDReq{})
	JSON(ctx, id, err)
}

func (c *APIContext) get(ctx *gin.Context) {
	ctx.Redirect(http.StatusMovedPermanently, ipfsGetURL("api/v0/get"))
}

func ipfsGetURL(uri string) string {
	return fmt.Sprintf("%s/%s", config.IPFSAddrHTTP(), uri)
}

func (c *APIContext) ping(ctx *gin.Context) {
	ping, err := c.Ping(&core.PingReq{})
	JSON(ctx, ping, err)
}

func (c *APIContext) debug(ctx *gin.Context) {
	uri := ctx.Query("uri")
	ctx.Redirect(http.StatusFound, ipfsGetURL(uri))
}

func (c *APIContext) query(ctx *gin.Context) {
	var err error
	j := struct {
		No string
	}{}
	err = ctx.BindJSON(&j)
	if err != nil {
		JSON(ctx, "", fmt.Errorf("query failed(%w)", err))
		return
	}
	dTag, e := c.c.DTag()
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

func (c *APIContext) nodeLink() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req core.NodeLinkReq
		err := ctx.BindJSON(&req)
		if err != nil {
			JSON(ctx, nil, err)
			return
		}
		id, err := c.Link(&req)
		JSON(ctx, id, err)
	}

}

func (c *APIContext) nodeUnlink() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id, err := c.Unlink(&core.NodeUnlinkReq{})
		JSON(ctx, id, err)
	}
}

// Unlink ...
func (c *APIContext) Unlink(req *core.NodeUnlinkReq) (*core.NodeUnlinkResp, error) {
	return c.m.NodeAPI().Unlink(req)
}

func (c *APIContext) nodeList() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		list, err := c.List(&core.NodeListReq{})
		JSON(ctx, list, err)
	}
}

func (c *APIContext) datastorePinLs() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		list, err := c.PinLs(&core.DataStoreReq{})
		JSON(ctx, list, err)
	}
}

func (c *APIContext) add() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req core.AddReq
		err := ctx.BindJSON(&req)
		if err != nil {
			JSON(ctx, nil, err)
			return
		}
		add, err := c.Add(&req)
		if err != nil {
			JSON(ctx, nil, err)
			return
		}
		JSON(ctx, add, nil)
	}
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
