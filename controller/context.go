package controller

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	ma "github.com/multiformats/go-multiaddr"
	mnet "github.com/multiformats/go-multiaddr-net"
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

// APIContext ...
type APIContext struct {
	cfg        *config.Config
	eng        *gin.Engine
	listener   net.Listener
	serv       *http.Server
	ready      *atomic.Bool
	controller *Controller
	//ethNode    *nodeBinETH
	//ipfsNode   *nodeBinIPFS
	msg func(s string)
	//cb         func(tag core.RequestTag, v interface{}) error
	manager core.NodeManager
}

// Add ...
func (c *APIContext) Add(req *core.AddReq) (*core.AddResp, error) {
	var info core.DataInfoV1
	err := info.Unmarshal([]byte(req.JSNFO))
	if err != nil {
		return nil, err
	}
	return &core.AddResp{}, nil
}

var _ core.API = &APIContext{}

// NewContext ...
func NewContext(cfg *config.Config) *APIContext {
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	eng := gin.Default()
	return &APIContext{
		cfg:   cfg,
		eng:   eng,
		ready: atomic.NewBool(false),
		serv: &http.Server{
			Handler: eng,
		},
	}
}

// API ...
func (c *APIContext) API(manager core.NodeManager) core.API {
	c.manager = manager
	return c
}

// NodeAddrInfo ...
func (c *APIContext) NodeAddrInfo(req *core.AddrReq) (*core.AddrResp, error) {
	if req.ID == "" {
		return &core.AddrResp{}, nil
	}
	panic("implement me")
}

// Ping ...
func (c *APIContext) Ping(req *core.PingReq) (*core.PingResp, error) {
	return &core.PingResp{
		Data: "pong",
	}, nil
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
	log.Infow("result id", "id", c.cfg.Identity, "public key", pubString)
	ipfsID, err := c.controller.dataNode().ID(context.TODO())
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
	v0.POST("/node/nodeLink", c.nodeLink)
	v0.POST("/node/unlink", c.nodeUnlink)
	v0.POST("/node/list", c.nodeList)
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

func (c *APIContext) setController(controller *Controller) {
	c.controller = controller
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
	dTag, e := c.controller.infoNode().DTag()
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

func (c *APIContext) nodeLink(ctx *gin.Context) {
	id, err := c.Link(&core.NodeLinkReq{})
	JSON(ctx, id, err)
}

// Link ...
func (c *APIContext) Link(req *core.NodeLinkReq) (*core.NodeLinkResp, error) {
	for _, addr := range req.Addrs {

		dial, err := mnet.Dial(addr)
		if err != nil {
			continue
		}
		c.manager.Conn(dial)
		return &core.NodeLinkResp{}, nil
	}
	return &core.NodeLinkResp{}, errors.New("all request was failed")
}

func (c *APIContext) nodeUnlink(ctx *gin.Context) {
	id, err := c.Unlink(&core.NodeUnlinkReq{})
	JSON(ctx, id, err)
}

// Unlink ...
func (c *APIContext) Unlink(req *core.NodeUnlinkReq) (*core.NodeUnlinkResp, error) {
	panic("implement me")
}
func (c *APIContext) nodeList(ctx *gin.Context) {

}

// List ...
func (c *APIContext) List(req *core.NodeListReq) (*core.NodeListResp, error) {
	panic("implement me")
}

// NodeAPI ...
func (c *APIContext) NodeAPI() core.NodeAPI {
	return c
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
