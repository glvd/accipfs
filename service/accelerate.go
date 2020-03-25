package service

import (
	"context"
	"errors"
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
	id         *core.NodeInfo
	nodes      core.NodeStore
	peerNodes  *atomic.Int64
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
		peerNodes:  atomic.NewInt64(0),
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

	jobAcc, err := a.cron.AddJob("0 1/3 * * * *", a)
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
		fmt.Println(outputHead, "Accelerate", "syncing node", info.Name)

		err := Ping(info)
		if err != nil {
			a.nodes.Remove(info.Name)
			a.peerNodes.Add(-1)
			a.dummyNodes.Add(info)
			return true
		}
		nodeInfos, err := Peers(info)
		if err != nil {
			return true
		}

		for _, nodeInfo := range nodeInfos {
			if a.peerNodes.Load() > a.cfg.Limit {
				return false
			}
			result := new(bool)
			if err := a.addPeer(ctx, nodeInfo, result); err != nil {
				continue
			}
		}
		//time.Sleep(30 * time.Second)
		return true
	})
	fmt.Println(outputHead, "Accelerate", "syncing done")
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

func (a *Accelerate) localID() (*core.NodeInfo, error) {
	var info core.NodeInfo
	info.Name = a.self.Name
	info.Version = core.Version
	info.Port = a.cfg.Port
	log.Debugw("print remote ip", "tag", outputHead, "ip", info.RemoteAddr, "port", info.Port)
	ds, err := a.ipfsClient.ID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("datastore error:%w", err)
	}
	info.DataStore = *ds
	c, err := a.ethClient.NodeInfo(context.Background())
	if err != nil {
		return nil, fmt.Errorf("nodeinfo error:%w", err)
	}
	info.Contract = *c
	return &info, nil
}

// ID ...
func (a *Accelerate) ID(r *http.Request, e *Empty, result *core.NodeInfo) error {
	id, err := a.localID()
	if err != nil {
		return err
	}
	*result = *id
	return nil
}

// Connect ...
func (a *Accelerate) Connect(r *http.Request, node *core.NodeInfo, result *core.NodeInfo) error {
	log.Infow("connect", "tag", outputHead, "addr", r.RemoteAddr)
	if node == nil {
		return fmt.Errorf("nil node info")
	}

	node.RemoteAddr, _ = general.SplitIP(r.RemoteAddr)

	id, err := a.localID()
	if err != nil {
		return err
	}
	*result = *id

	err = Ping(node)
	if err != nil {
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

func (a *Accelerate) addPeer(ctx context.Context, info *core.NodeInfo, result *bool) error {
	*result = false

	//skip self add
	if info.Name == a.id.Name {
		return errors.New("cannot add your self")
	}

	err := Ping(info)
	if err != nil {
		log.Errorw("add peer", "tag", outputHead, "error", err)
		a.dummyNodes.Add(info)
		return nil
	}

	ipfsTimeout, cancelFunc := context.WithTimeout(ctx, time.Duration(a.cfg.Interval)*time.Second)
	var ipfsErr error
	for _, addr := range info.DataStore.Addresses {
		ipfsErr = a.ipfsClient.SwarmConnect(ipfsTimeout, addr)
		if ipfsErr == nil {
			break
		}
	}
	cancelFunc()
	if ipfsErr != nil {
		a.dummyNodes.Add(info)
		log.Errorw("add peer", "tag", outputHead, "error", ipfsErr)

		return nil
	}
	ethTimeout, cancelFunc := context.WithTimeout(ctx, time.Duration(a.cfg.Interval)*time.Second)
	//fmt.Println("connect eth:", info.Contract.Enode)
	err = a.ethClient.AddPeer(ethTimeout, info.Contract.Enode)
	if err != nil {
		a.dummyNodes.Add(info)
		log.Errorw("add peer", "tag", outputHead, "error", err)
		return nil
	}
	cancelFunc()

	a.peerNodes.Add(1)
	a.nodes.Add(info)
	*result = true
	return nil
}

// AddPeer ...
func (a *Accelerate) AddPeer(r *http.Request, info *core.NodeInfo, result *bool) error {
	return a.addPeer(r.Context(), info, result)
}

// Peers ...
func (a *Accelerate) Peers(r *http.Request, empty *Empty, result *[]*core.NodeInfo) error {
	a.nodes.Range(func(info *core.NodeInfo) bool {
		*result = append(*result, info)
		return true
	})
	return nil
}

// Exchange ...
func (a *Accelerate) Exchange(r *http.Request, from, to interface{}) error {
	return nil
}
