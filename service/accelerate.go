package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/glvd/accipfs/client"
	"github.com/glvd/accipfs/controller"
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

// BustLinker ...
type BustLinker struct {
	id         *core.NodeInfo
	tasks      task.Task
	cache      *cache.MemoryCache
	nodes      core.NodeStore
	dummyNodes core.NodeStore
	lock       *atomic.Bool
	self       *account.Account
	cfg        *config.Config
	c          *controller.Controller
	cron       *cron.Cron
}

// BootList ...
var BootList = []string{
	"gate.dhash.app",
}

// NewBustLinker ...
func NewBustLinker(cfg *config.Config) (linker *BustLinker, err error) {
	linker = &BustLinker{
		nodes:      core.NewNodeStore(),
		dummyNodes: core.NewNodeStore(),
		lock:       atomic.NewBool(false),
		cfg:        cfg,
	}
	//linker.ethServer = newNodeServerETH(cfg)
	//linker.ipfsServer = newNodeServerIPFS(cfg)
	//linker.ethClient, _ = newNodeETH(cfg)
	//linker.ipfsClient, _ = newNodeIPFS(cfg)
	linker.cache = cache.New(cfg)
	linker.tasks = task.New()
	linker.cron = cron.New(cron.WithSeconds())
	selfAcc, err := account.LoadAccount(cfg)
	if err != nil {
		return nil, err
	}

	linker.self = selfAcc
	return linker, nil
}

// Start ...
func (l *BustLinker) Start() {
	go l.c.Run()

	jobAcc, err := l.cron.AddJob("0 1/3 * * * *", l)
	if err != nil {
		panic(err)
	}
	output("BustLinker", "run id", jobAcc)
	l.cron.Run()
}

// Run ...
func (l *BustLinker) Run() {
	if l.lock.Load() {
		output("BustLinker", "the previous task has not been completed")
		return
	}
	l.lock.Store(true)
	defer l.lock.Store(false)
	ctx := context.TODO()
	l.nodes.Range(func(info *core.NodeInfo) bool {
		output("BustLinker", "syncing node", info.Name)

		err := client.Ping(info)
		if err != nil {
			l.nodes.Remove(info.Name)
			l.dummyNodes.Add(info)
			logE("ping failed", "account", info.Name, "error", err)
			return true
		}
		url := info.Address().URL()
		nodeInfos, err := client.Peers(url, info)
		if err != nil {
			logE("get peers failed", "account", info.Name, "error", err)
			return true
		}

		for _, nodeInfo := range nodeInfos {
			if l.nodes.Length() > l.cfg.Limit {
				return false
			}
			result := new(bool)
			if err := l.addPeer(ctx, nodeInfo, result); err != nil {
				logE("add peer failed", "account", info.Name, "error", err)
				continue
			}
			if *result {
				pins, err := client.Pins(nodeInfo)
				if err != nil {
					logE("get pin list", "error", err)
					continue
				}
				for _, p := range pins {
					err := l.cache.AddOrUpdate(p, nodeInfo)
					if err != nil {
						logE("cache add or update", "error", err)
						continue
					}
				}
			}

		}
		//time.Sleep(30 * time.Second)
		return true
	})
	fmt.Println(outputHead, "BustLinker", "syncing done")
}

// Stop ...
func (l *BustLinker) Stop() {
	ctx := l.cron.Stop()
	<-ctx.Done()
	if err := l.ethServer.Stop(); err != nil {
		log.Errorw("eth stop error", "tag", outputHead, "error", err)
		return
	}

	if err := l.ipfsServer.Stop(); err != nil {
		log.Errorw("ipfs stop error", "tag", outputHead, "error", err)
		return
	}

}

// Ping ...
func (l *BustLinker) Ping(r *http.Request, e *core.Empty, result *string) error {
	*result = "pong"
	return nil
}

func (l *BustLinker) localID() (*core.NodeInfo, error) {
	var info core.NodeInfo
	info.Name = l.self.Name
	info.Version = core.Version
	info.RemoteAddr = "127.0.0.1"
	info.Port = l.cfg.Port
	log.Debugw("print remote ip", "tag", outputHead, "ip", info.RemoteAddr, "port", info.Port)
	ds, err := l.ipfsClient.ID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("datastore error:%w", err)
	}
	info.DataStore = *ds
	c, err := l.ethClient.NodeInfo(context.Background())
	if err != nil {
		return nil, fmt.Errorf("nodeinfo error:%w", err)
	}
	info.Contract = *c
	return &info, nil
}

// ID ...
func (l *BustLinker) ID(r *http.Request, e *core.Empty, result *core.NodeInfo) error {
	id, err := l.localID()
	if err != nil {
		return err
	}
	*result = *id
	return nil
}

// Connected ...
func (l *BustLinker) Connected(r *http.Request, node *core.NodeInfo, result *core.NodeInfo) error {
	log.Infow("connected", "tag", outputHead, "addr", r.RemoteAddr)
	if node == nil {
		return fmt.Errorf("nil node info")
	}

	node.RemoteAddr, _ = general.SplitIP(r.RemoteAddr)

	id, err := l.localID()
	if err != nil {
		return err
	}
	*result = *id

	err = client.Ping(node)
	if err != nil {
		if !l.dummyNodes.Check(node.Name) {
			l.dummyNodes.Add(node)
		}
		return nil
	}
	if !l.nodes.Check(node.Name) {
		l.nodes.Add(node)
		return nil
	}
	return nil
}

// ConnectTo ...
func (l BustLinker) ConnectTo(r *http.Request, addr *string, result *core.NodeInfo) error {
	id, err := l.localID()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("http://%s/rpc", *addr)

	err = general.RPCPost(url, "BustLinker.Connected", id, result)
	if err != nil {
		return err
	}
	result.RemoteAddr, result.Port = general.SplitIP(*addr)
	return nil
}

func (l *BustLinker) addPeer(ctx context.Context, info *core.NodeInfo, result *bool) error {
	*result = false

	if info.Name == l.id.Name {
		//ignore self add
		return nil
	}

	err := client.Ping(info)
	if err != nil {
		log.Errorw("add peer", "tag", outputHead, "error", err)
		l.dummyNodes.Add(info)
		return err
	}

	ipfsTimeout, cancelFunc := context.WithTimeout(ctx, time.Duration(l.cfg.Interval)*time.Second)
	var ipfsErr error
	for _, addr := range info.DataStore.Addresses {
		ipfsErr = l.ipfsClient.SwarmConnect(ipfsTimeout, addr)
		if ipfsErr == nil {
			break
		}
	}
	cancelFunc()
	if ipfsErr != nil {
		l.dummyNodes.Add(info)
		log.Errorw("add peer", "tag", outputHead, "error", ipfsErr)

		return err
	}
	ethTimeout, cancelFunc := context.WithTimeout(ctx, time.Duration(l.cfg.Interval)*time.Second)
	//fmt.Println("connect eth:", info.Contract.Enode)
	err = l.ethClient.AddPeer(ethTimeout, info.Contract.Enode)
	if err != nil {
		l.dummyNodes.Add(info)
		log.Errorw("add peer", "tag", outputHead, "error", err)
		return err
	}
	cancelFunc()

	l.nodes.Add(info)
	*result = true
	return nil
}

// AddPeer ...
func (l *BustLinker) AddPeer(r *http.Request, info *core.NodeInfo, result *bool) error {
	return l.addPeer(r.Context(), info, result)
}

// Peers ...
func (l *BustLinker) Peers(r *http.Request, _ *core.Empty, result *[]*core.NodeInfo) error {
	l.nodes.Range(func(info *core.NodeInfo) bool {
		*result = append(*result, info)
		return true
	})
	return nil
}

func (l *BustLinker) pins(ctx context.Context, result *[]string) error {
	pins, e := l.ipfsClient.PinLS(ctx)
	if e != nil {
		return e
	}
	for _, p := range pins {
		*result = append(*result, p.Path().String())
	}
	return nil
}

// Pins ...
func (l *BustLinker) Pins(r *http.Request, _ *core.Empty, result *[]string) error {
	return l.pins(r.Context(), result)
}

// PinVideo ...
func (l *BustLinker) PinVideo(r *http.Request, no *string, result *bool) error {
	info := new(string)
	err := l.tagInfo(*no, info)
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
		err := l.nodeConnect(ctx, v.PosterHash)
		if err != nil {
			cancelFunc()
			resultErr <- err
			return
		}
		e := l.ipfsClient.PinAdd(ctx, v.PosterHash)
		if e != nil {
			cancelFunc()
			resultErr <- e
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := l.nodeConnect(ctx, v.ThumbHash)
		if err != nil {
			cancelFunc()
			resultErr <- err
			return
		}
		e := l.ipfsClient.PinAdd(ctx, v.ThumbHash)
		if e != nil {
			cancelFunc()
			resultErr <- e
		}

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := l.nodeConnect(ctx, v.SourceHash)
		if err != nil {
			cancelFunc()
			resultErr <- err
			return
		}
		e := l.ipfsClient.PinAdd(ctx, v.SourceHash)
		if e != nil {
			cancelFunc()
			resultErr <- e
		}

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := l.nodeConnect(ctx, v.M3U8Hash)
		if err != nil {
			cancelFunc()
			resultErr <- err
			return
		}
		e := l.ipfsClient.PinAdd(ctx, v.M3U8Hash)
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

func (l *BustLinker) tagInfo(tag string, info *string) error {
	dTag, e := l.ethClient.DTag()
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
func (l *BustLinker) TagInfo(_ *http.Request, tag *string, info *string) error {
	return l.tagInfo(*tag, info)
}

// Info ...
func (l *BustLinker) Info(r *http.Request, hash *string, info *string) error {
	bytes, e := l.cache.Get(*hash)
	if e != nil {
		return e
	}
	*info = string(bytes)
	return nil
}

// Exchange ...
func (l *BustLinker) Exchange(r *http.Request, n *core.NodeInfo, to []string) error {

	return nil
}

func (l *BustLinker) nodeConnect(ctx context.Context, hash string) error {
	hashInfo, err := l.cache.GetHashInfo(hash)
	if err != nil {
		return err
	}
	for info := range hashInfo {
		nodeInfo, err := l.cache.GetNodeInfo(info)
		if err != nil {
			continue
		}
		var resultErr error
		for _, addr := range nodeInfo.DataStore.Addresses {
			resultErr = l.ipfsClient.SwarmConnect(ctx, addr)
			if resultErr == nil {
				break
			}
		}
	}
	return nil
}
