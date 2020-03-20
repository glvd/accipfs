package service

import (
	"context"
	"fmt"
	"github.com/glvd/accipfs/account"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/glvd/accipfs/general"
	"github.com/goextension/log"
	"github.com/robfig/cron/v3"
	"go.uber.org/atomic"
	"net/http"
	"time"
)

// Empty ...
type Empty struct {
}

// Accelerate ...
type Accelerate struct {
	nodes      core.NodeStore
	dummyNodes core.NodeStore
	lock       *atomic.Bool
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
		lock:       atomic.NewBool(false),
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

// Start ...
func (a *Accelerate) Start() {
	if err := a.ethServer.Start(); err != nil {
		panic(err)
	}
	if err := a.ipfsServer.Start(); err != nil {
		panic(err)
	}

	//ethNode, err := a.ethServer.Node()
	//jobETH, err := a.cron.AddJob("0 * * * * *", ethNode)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(outputHead, "ETH", "run id", jobETH)

	//ipfsNode, err := a.ipfsServer.Node()
	//jobIPFS, err := a.cron.AddJob("0 * * * * *", ipfsNode)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(outputHead, "IPFS", "run id", jobIPFS)

	jobAcc, err := a.cron.AddJob("0 0/30 * * *", a)
	if err != nil {
		panic(err)
	}
	fmt.Println(outputHead, "Accelerate", "run id", jobAcc)
	a.cron.Run()
}

// Run ...
func (a *Accelerate) Run() {
	if a.lock.Load() {
		fmt.Println(outputHead, "Accelerate", "accelerate.run is already running")
		return
	}
	a.lock.Store(true)
	defer a.lock.Store(false)
	ctx := context.TODO()

	a.nodes.Range(func(info *core.NodeInfo) bool {
		err := Ping(info)
		if err != nil {
			a.nodes.Remove(info.Name)
			a.dummyNodes.Add(info)
			return true
		}
		nodeInfos, err := Peers(info)
		if err != nil {
			return true
		}

		for _, nodeInfo := range nodeInfos {
			err := Ping(nodeInfo)
			if err != nil {
				a.dummyNodes.Add(nodeInfo)
				continue
			}

			ipfsTimeout, cancelFunc := context.WithTimeout(ctx, time.Duration(a.cfg.Interval)*time.Second)
			var ipfsErr error
			for _, addr := range nodeInfo.DataStore.Addresses {
				ipfsErr = a.ipfsClient.SwarmConnect(ipfsTimeout, addr)
				if ipfsErr == nil {
					break
				}
			}
			cancelFunc()
			if ipfsErr != nil {
				a.dummyNodes.Add(nodeInfo)
				continue
			}
			ethTimeout, cancelFunc := context.WithTimeout(ctx, time.Duration(a.cfg.Interval)*time.Second)
			err = a.ethClient.AddPeer(ethTimeout, nodeInfo.Contract.Enode)
			if err != nil {
				a.dummyNodes.Add(nodeInfo)
				continue
			}
			cancelFunc()
			a.nodes.Add(nodeInfo)
		}
		time.Sleep(3 * time.Second)
		return true
	})
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

// ID ...
func (a *Accelerate) ID(r *http.Request, e *Empty, result *core.NodeInfo) error {
	result.Name = a.self.Name
	result.Version = core.Version
	result.Port = a.cfg.Port
	fmt.Println(outputHead, "Accelerate", "print remote ip:", result.RemoteAddr, ":", result.Port)
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

	err := Ping(node)
	if err != nil {
		*result = false
		if !a.dummyNodes.Check(node.Name) {
			a.dummyNodes.Add(node)
		}
		return nil
	}
	if !a.nodes.Check(node.Name) {
		a.nodes.Add(node)
		return nil
	}
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
