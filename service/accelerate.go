package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/glvd/accipfs/client"
	"github.com/glvd/accipfs/task"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/glvd/accipfs/account"
	"github.com/glvd/accipfs/cache"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/glvd/accipfs/general"
	"github.com/goextension/log"
	"github.com/robfig/cron/v3"
	"go.uber.org/atomic"
)

// Accelerate ...
type Accelerate struct {
	id         *core.NodeInfo
	tasks      task.Task
	cache      *cache.MemoryCache
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
	acc.cache = cache.New(cfg)
	acc.tasks = task.New()
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
		fmt.Println(outputHead, "Accelerate", "the previous task has not been completed")
		return
	}
	a.lock.Store(true)
	defer a.lock.Store(false)
	ctx := context.TODO()
	a.nodes.Range(func(info *core.NodeInfo) bool {
		fmt.Println(outputHead, "Accelerate", "syncing node", info.Name)

		err := client.Ping(info)
		if err != nil {
			a.nodes.Remove(info.Name)
			a.dummyNodes.Add(info)
			return true
		}
		url := info.Address().URL()
		nodeInfos, err := client.Peers(url, info)
		if err != nil {
			return true
		}

		for _, nodeInfo := range nodeInfos {
			if a.nodes.Length() > a.cfg.Limit {
				return false
			}
			result := new(bool)
			if err := a.addPeer(ctx, nodeInfo, result); err != nil {
				continue
			}
			if *result {
				pins, err := client.Pins(nodeInfo)
				if err != nil {
					log.Errorw("get pin list", "error", err)
					continue
				}
				for _, p := range pins {
					err := a.cache.AddOrUpdate(p, nodeInfo)
					if err != nil {
						log.Errorw("cache add or update", "error", err)
						continue
					}
				}
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
func (a *Accelerate) Ping(r *http.Request, e *core.Empty, result *string) error {
	*result = "pong"
	return nil
}

func (a *Accelerate) localID() (*core.NodeInfo, error) {
	var info core.NodeInfo
	info.Name = a.self.Name
	info.Version = core.Version
	info.RemoteAddr = "127.0.0.1"
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
func (a *Accelerate) ID(r *http.Request, e *core.Empty, result *core.NodeInfo) error {
	id, err := a.localID()
	if err != nil {
		return err
	}
	*result = *id
	return nil
}

// Connected ...
func (a *Accelerate) Connected(r *http.Request, node *core.NodeInfo, result *core.NodeInfo) error {
	log.Infow("connected", "tag", outputHead, "addr", r.RemoteAddr)
	if node == nil {
		return fmt.Errorf("nil node info")
	}

	node.RemoteAddr, _ = general.SplitIP(r.RemoteAddr)

	id, err := a.localID()
	if err != nil {
		return err
	}
	*result = *id

	err = client.Ping(node)
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

// ConnectTo ...
func (a Accelerate) ConnectTo(r *http.Request, addr *string, result *core.NodeInfo) error {
	id, err := a.localID()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("http://%s/rpc", *addr)

	err = general.RPCPost(url, "Accelerate.Connected", id, result)
	if err != nil {
		return err
	}
	result.RemoteAddr, result.Port = general.SplitIP(*addr)
	return nil
}

func (a *Accelerate) addPeer(ctx context.Context, info *core.NodeInfo, result *bool) error {
	*result = false

	if info.Name == a.id.Name {
		//ignore self add
		return nil
	}

	err := client.Ping(info)
	if err != nil {
		log.Errorw("add peer", "tag", outputHead, "error", err)
		a.dummyNodes.Add(info)
		return err
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

		return err
	}
	ethTimeout, cancelFunc := context.WithTimeout(ctx, time.Duration(a.cfg.Interval)*time.Second)
	//fmt.Println("connect eth:", info.Contract.Enode)
	err = a.ethClient.AddPeer(ethTimeout, info.Contract.Enode)
	if err != nil {
		a.dummyNodes.Add(info)
		log.Errorw("add peer", "tag", outputHead, "error", err)
		return err
	}
	cancelFunc()

	a.nodes.Add(info)
	*result = true
	return nil
}

// AddPeer ...
func (a *Accelerate) AddPeer(r *http.Request, info *core.NodeInfo, result *bool) error {
	return a.addPeer(r.Context(), info, result)
}

// Peers ...
func (a *Accelerate) Peers(r *http.Request, _ *core.Empty, result *[]*core.NodeInfo) error {
	a.nodes.Range(func(info *core.NodeInfo) bool {
		*result = append(*result, info)
		return true
	})
	return nil
}

func (a *Accelerate) pins(ctx context.Context, result *[]string) error {
	pins, e := a.ipfsClient.PinLS(ctx)
	if e != nil {
		return e
	}
	for _, p := range pins {
		*result = append(*result, p.Path().String())
	}
	return nil
}

// Pins ...
func (a *Accelerate) Pins(r *http.Request, _ *core.Empty, result *[]string) error {
	return a.pins(r.Context(), result)
}

// PinVideo ...
func (a *Accelerate) PinVideo(r *http.Request, no *string, result *bool) error {
	info := new(string)
	err := a.tagInfo(*no, info)
	if err != nil {
		return err
	}

	var v core.VideoV1
	reader := strings.NewReader(*info)
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&v)
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	resultErr := make(chan error)
	ctx, cancelFunc := context.WithCancel(r.Context())
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := a.nodeConnect(ctx, v.PosterHash)
		if err != nil {
			cancelFunc()
			resultErr <- err
			return
		}
		e := a.ipfsClient.PinAdd(ctx, v.PosterHash)
		if e != nil {
			cancelFunc()
			resultErr <- e
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := a.nodeConnect(ctx, v.ThumbHash)
		if err != nil {
			cancelFunc()
			resultErr <- err
			return
		}
		e := a.ipfsClient.PinAdd(ctx, v.ThumbHash)
		if e != nil {
			cancelFunc()
			resultErr <- e
		}

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := a.nodeConnect(ctx, v.SourceHash)
		if err != nil {
			cancelFunc()
			resultErr <- err
			return
		}
		e := a.ipfsClient.PinAdd(ctx, v.SourceHash)
		if e != nil {
			cancelFunc()
			resultErr <- e
		}

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := a.nodeConnect(ctx, v.M3U8Hash)
		if err != nil {
			cancelFunc()
			resultErr <- err
			return
		}
		e := a.ipfsClient.PinAdd(ctx, v.M3U8Hash)
		if e != nil {
			cancelFunc()
			resultErr <- e
		}
	}()

	wg.Wait()
	select {
	case e := <-resultErr:
		return e
	default:
	}
	*result = true
	return nil
}

func (a *Accelerate) tagInfo(tag string, info *string) error {
	dTag, e := a.ethClient.DTag()
	if e != nil {
		return e
	}
	message, e := dTag.GetTagMessage(&bind.CallOpts{Pending: true}, "video", tag)
	if e != nil {
		return e
	}

	if message.Size.Int64() > 0 {
		*info = message.Value[0]
	}
	return nil
}

// TagInfo ...
func (a *Accelerate) TagInfo(_ *http.Request, tag *string, info *string) error {
	return a.tagInfo(*tag, info)
}

// Info ...
func (a *Accelerate) Info(r *http.Request, hash *string, info *string) error {
	bytes, e := a.cache.Get(*hash)
	if e != nil {
		return e
	}
	*info = string(bytes)
	return nil
}

// Exchange ...
func (a *Accelerate) Exchange(r *http.Request, n *core.NodeInfo, to []string) error {

	return nil
}

func (a *Accelerate) nodeConnect(ctx context.Context, hash string) error {
	hashInfo, err := a.cache.GetHashInfo(hash)
	if err != nil {
		return err
	}
	for info := range hashInfo {
		nodeInfo, err := a.cache.GetNodeInfo(info)
		if err != nil {
			continue
		}
		var resultErr error
		for _, addr := range nodeInfo.DataStore.Addresses {
			resultErr = a.ipfsClient.SwarmConnect(ctx, addr)
			if resultErr == nil {
				break
			}
		}
	}
	return nil
}
