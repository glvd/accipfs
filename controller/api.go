package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
	ethNode    *nodeBinETH
	ipfsNode   *nodeBinIPFS
}

// Ping ...
func (a *API) Ping(req *core.PingReq) (*core.PingResp, error) {
	panic("implement me")
}

// ID ...
func (a *API) ID(req *core.IDReq) (*core.IDResp, error) {
	ctx := context.Background()
	id, err := a.ipfsNode.ID(ctx)
	if err != nil {
		return nil, err
	}
	info, err := a.ethNode.NodeInfo(ctx)
	if err != nil {
		return nil, err
	}
	return &core.IDResp{
		Name:      "",
		DataStore: id,
		Contract:  info,
	}, nil
}

// New ...
func newAPI(cfg *config.Config, controller *Controller) core.API {
	eng := gin.Default()
	return &API{
		cfg:      cfg,
		ethNode:  controller.services[IndexETH].(*nodeBinETH),
		ipfsNode: controller.services[IndexIPFS].(*nodeBinIPFS),
		eng:      eng,
		serv: &http.Server{
			Handler: eng,
		},
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

func (a *API) routeList() {
	api := a.eng.Group("/api")
	api.GET("/ping", a.ping)
	g := api.Group(a.cfg.API.Version)
	g.Handle(http.MethodGet, "/id", a.id)

	if a.cfg.Debug {
		g.GET("/debug", a.debug)
	}

	v0 := g.Group("v0")
	v0.GET("/get", a.get)
	v0.GET("/query", a.query)
}

// Stop ...
func (a *API) Stop() error {
	if a.serv != nil {
		return a.serv.Close()
	}
	return nil
}

func (a *API) id(c *gin.Context) {
	id, err := a.ID(&core.IDReq{})
	JSON(c, id, err)
}

func (a *API) get(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, ipfsGetURL("api/v0/get"))
}

func ipfsGetURL(uri string) string {
	return fmt.Sprintf("%s/%s", config.IPFSAddrHTTP(), uri)
}

func (a *API) ping(c *gin.Context) {
	ping, err := a.Ping(&core.PingReq{})
	JSON(c, ping, err)
}

func (a *API) debug(c *gin.Context) {
	uri := c.Query("uri")
	c.Redirect(http.StatusMovedPermanently, ipfsGetURL(uri))
}

func (a *API) query(c *gin.Context) {
	var err error
	j := struct {
		No string
	}{}
	err = c.BindJSON(&j)
	if err != nil {
		JSON(c, "", fmt.Errorf("query failed(%w)", err))
		return
	}
	dTag, e := a.ethNode.DTag()
	if e != nil {
		JSON(c, "", fmt.Errorf("query failed(%w)", e))
		return
	}
	message, e := dTag.GetTagMessage(&bind.CallOpts{Pending: true}, "video", j.No)
	if e != nil {
		JSON(c, "", fmt.Errorf("query failed(%w)", e))
		return
	}

	if message.Size.Int64() > 0 {
		JSON(c, message.Value[0], nil)
		return
	}
	JSON(c, "", nil)
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
		"message": m,
	})

}
