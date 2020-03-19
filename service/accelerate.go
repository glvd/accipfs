package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/glvd/accipfs/account"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/glvd/accipfs/general"
	"github.com/goextension/log"
	"github.com/gorilla/rpc/v2/json2"
	"github.com/robfig/cron/v3"
	"net/http"
)

// Empty ...
type Empty struct {
}

// Accelerate ...
type Accelerate struct {
	nodes      core.NodeStore
	dummyNodes core.NodeStore
	self       *account.Account
	cfg        *config.Config
	ethServer  *nodeServerETH
	ethClient  *nodeClientETH
	ipfsServer *nodeServerIPFS
	ipfsClient *nodeClientIPFS
	cron       *cron.Cron
}

// BootList ...
var BootList = []string{
	"gate.dhash.app",
}

// NewAccelerateServer ...
func NewAccelerateServer(cfg *config.Config) (acc *Accelerate, err error) {
	acc = &Accelerate{
		nodes:      core.NewNodeStore(),
		dummyNodes: core.NewNodeStore(),
		cfg:        cfg,
	}
	acc.ethServer = newNodeServerETH(cfg)
	acc.ipfsServer = newNodeServerIPFS(cfg)
	acc.ethClient, _ = newNodeETH(cfg)
	acc.ipfsClient, _ = newNodeIPFS(cfg)
	acc.cron = cron.New(cron.WithSeconds())
	selfAcc, err := account.LoadAccount(cfg)
	if err != nil {
		return nil, err
	}

	acc.self = selfAcc
	return acc, nil
}

// Run ...
func (a *Accelerate) Run() {
	fmt.Println("accelerate run")
	if err := a.ethServer.Start(); err != nil {
		panic(err)
	}
	if err := a.ipfsServer.Start(); err != nil {
		panic(err)
	}

	ethNode, err := a.ethServer.Node()
	jobETH, err := a.cron.AddJob("0 * * * * *", ethNode)
	if err != nil {
		panic(err)
	}
	fmt.Println(outputHead, "ETH", "run id", jobETH)

	ipfsNode, err := a.ipfsServer.Node()
	jobIPFS, err := a.cron.AddJob("0 * * * * *", ipfsNode)
	if err != nil {
		panic(err)
	}
	fmt.Println(outputHead, "IPFS", "run id", jobIPFS)

	a.cron.Run()
}

// Stop ...
func (a *Accelerate) Stop() {
	ctx := a.cron.Stop()
	<-ctx.Done()
	if err := a.ethServer.Stop(); err != nil {
		log.Errorw("eth stop error", "tag", outputHead, "error", err)
		return
	}

	if err := a.ipfsServer.Stop(); err != nil {
		log.Errorw("ipfs stop error", "tag", outputHead, "error", err)
		return
	}
}

// Ping ...
func (a *Accelerate) Ping(r *http.Request, e *Empty, result *string) error {
	*result = "pong"
	return nil
}

// Ping ...
func Ping(url string) error {
	pingReq, err := json2.EncodeClientRequest("Accelerate.Ping", &Empty{})
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(pingReq))
	if err != nil {
		return err
	}
	result := new(string)
	err = json2.DecodeClientResponse(resp.Body, result)
	if err != nil {
		return err
	}
	if *result != "pong" {
		return fmt.Errorf("get wrong response data:%s", *result)
	}

}

// ID ...
func (a *Accelerate) ID(r *http.Request, e *Empty, result *core.NodeInfo) error {
	result.Name = a.self.Name
	result.Version = core.Version
	fmt.Println(outputHead, "Accelerate", "print remote ip:", r.RemoteAddr)
	ds, err := a.ipfsClient.ID(context.Background())
	if err != nil {
		return fmt.Errorf("datastore error:%w", err)
	}
	result.DataStore = *ds
	c, err := a.ethClient.NodeInfo(context.Background())
	if err != nil {
		return fmt.Errorf("nodeinfo error:%w", err)
	}
	result.Contract = *c
	return nil
}

// Connect ...
func (a *Accelerate) Connect(r *http.Request, node *core.NodeInfo, result *bool) error {
	log.Infow("connect", "tag", outputHead, "addr", r.RemoteAddr)
	if node == nil {
		return fmt.Errorf("nil node info")
	}
	*result = true
	node.RemoteAddr, _ = general.SplitIP(r.RemoteAddr)

	if a.nodes.Check(node.Name) {
		*result = false
		return nil
	}
	a.nodes.Add(node)
	return nil
}

// Peers ...
func (a *Accelerate) Peers(r *http.Request, e *Empty, result []*core.NodeInfo) error {
	a.nodes.Range(func(info *core.NodeInfo) bool {
		result = append(result, info)
		return true
	})
	return nil
}

// Exchange ...
func (a *Accelerate) Exchange(r *http.Request, from, to interface{}) error {
	return nil
}
