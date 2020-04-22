package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/glvd/accipfs/account"
	"github.com/glvd/accipfs/cache"
	"github.com/glvd/accipfs/client"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/glvd/accipfs/general"
	"github.com/glvd/accipfs/task"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/robfig/cron/v3"
	"go.uber.org/atomic"
)

// BustLinker ...
type BustLinker struct {
	id    *core.Node
	tasks task.Task
	cache *cache.MemoryCache
	nodes NodeManager
	lock  *atomic.Bool
	self  *account.Account
	cfg   *config.Config
	eth   *ethNode
	ipfs  *ipfsNode
	cron  *cron.Cron
}

// BootList ...
var BootList = []string{
	"gate.dhash.app",
}

// NewBustLinker ...
func NewBustLinker(cfg *config.Config) (linker *BustLinker, err error) {
	linker = &BustLinker{
		nodes: NewNodeManager(),
		lock:  atomic.NewBool(false),
		cfg:   cfg,
	}
	//linker.ethServer = newNodeServerETH(cfg)
	//linker.ipfsServer = newNodeServerIPFS(cfg)
	linker.eth, _ = newNodeETH(cfg)
	linker.ipfs, _ = newNodeIPFS(cfg)
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
	//jobAcc, err := l.cron.AddJob("0 1/3 * * * *", l)
	jobAcc, err := l.cron.AddJob("0/5 * * * * *", l)
	if err != nil {
		panic(err)
	}
	output("bust linker", "run id", jobAcc)
	go l.cron.Run()
}

// Run ...
func (l *BustLinker) Run() {
	if l.lock.Load() {
		output("bust linker", "the previous task has not been completed")
		return
	}
	l.lock.Store(true)
	defer l.lock.Store(false)
	ctx := context.TODO()
	l.nodes.Range(func(node *core.Node) bool {
		output("bust linker", "syncing node", node.Name)
		v := l.nodes.Validate(node.NodeInfo.Name, func(node *core.Node) bool {
			err := client.Ping(general.RPCAddress(node.NodeAddress))
			if err != nil {
				logE("ping failed", "account", node.Name, "error", err)
				return false
			}
			return true
		})

		output("bust linker", "get peers", node.Name)
		if v {
			remoteNodes, err := client.Peers(general.RPCAddress(node.NodeAddress), node)
			if err != nil {
				logE("get peers failed", "account", node.Name, "error", err)
				return true
			}

			for _, rnode := range remoteNodes {
				if l.nodes.Length() > l.cfg.Limit {
					return false
				}
				result := new(bool)
				output("bust linker", "add peer", rnode.Name)
				if err := l.addPeer(ctx, rnode, result); err != nil {
					logE("add peer failed", "account", rnode.Name, "error", err)
					continue
				}
				if *result {
					output("bust linker", "pin source ", rnode.Name)
					pins, err := client.Pins(rnode.NodeAddress)
					if err != nil {
						logE("get pin list", "error", err)
						continue
					}
					for _, p := range pins {
						err := l.cache.AddOrUpdate(p, &rnode.NodeInfo)
						if err != nil {
							logE("cache add or update", "error", err)
							continue
						}
					}
				}

			}
		}
		return true
	})
	output("bust linker", "syncing done")
}

// WaitingForReady ...
func (l *BustLinker) WaitingForReady() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if l.eth.IsReady() {
			return
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if l.ipfs.IsReady() {
			return
		}
	}()
	wg.Wait()
}

// Stop ...
func (l *BustLinker) Stop() {
	ctx := l.cron.Stop()
	<-ctx.Done()
}

// Ping ...
func (l *BustLinker) Ping(r *http.Request, req *core.PingReq, resp *core.PingResp) error {
	resp.Resp = "pong"
	return nil
}

func (l *BustLinker) localID() (*core.Node, error) {
	var info core.Node
	info.Name = l.self.Name
	info.ProtocolVersion = core.Version
	info.Address = "127.0.0.1"
	info.Port = l.cfg.Port
	logD("print remote ip", "ip", info.Address, "port", info.Port)
	ds, err := l.ipfs.ID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("datastore error:%w", err)
	}
	info.DataStore = *ds
	c, err := l.eth.NodeInfo(context.Background())
	if err != nil {
		return nil, fmt.Errorf("nodeinfo error:%w", err)
	}
	info.Contract = *c
	return &info, nil
}

// ID ...
func (l *BustLinker) ID(r *http.Request, req *core.IDReq, resp *core.IDResp) error {
	id, err := l.localID()
	if err != nil {
		return err
	}
	resp.Node = *id
	return nil
}

// Connected ...
func (l *BustLinker) Connected(r *http.Request, req *core.ConnectedReq, resp *core.ConnectedResp) error {
	return l.connected(r, &req.Node, &resp.Node)
}

// Connected ...
func (l *BustLinker) connected(r *http.Request, node *core.Node, result *core.Node) error {
	logI("connected", "addr", r.RemoteAddr)
	if node == nil {
		return fmt.Errorf("nil node info")
	}

	node.NodeAddress.Address, _ = general.SplitIP(r.RemoteAddr)

	id, err := l.localID()
	if err != nil {
		return err
	}
	*result = *id
	err = client.Ping(general.RPCAddress(node.NodeAddress))
	if err != nil {
		return err
	}

	l.nodes.Add(node)
	return nil
}

// ConnectTo ...
func (l *BustLinker) ConnectTo(r *http.Request, req *core.ConnectToReq, resp *core.ConnectToResp) error {
	return l.connectTo(r, &req.Addr, &resp.Node)
}

// ConnectTo ...
func (l *BustLinker) connectTo(r *http.Request, addr *string, respNode *core.Node) error {
	id, err := l.localID()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("http://%s/rpc", *addr)
	connReq := &core.ConnectedReq{Node: *id}
	resp := new(core.ConnectedResp)
	err = general.RPCPost(url, "BustLinker.Connected", connReq, resp)
	if err != nil {
		return err
	}
	*respNode = resp.Node
	respNode.NodeAddress.Address, respNode.NodeAddress.Port = general.SplitIP(*addr)

	return nil
}

func (l *BustLinker) addPeer(ctx context.Context, node *core.Node, result *bool) error {
	*result = false

	if node.Name == l.id.Name {
		//ignore self add
		return nil
	}

	faultNode := l.nodes.IsFault(node.NodeInfo.Name)

	if faultNode != nil {
		if remain, b := faultTimeCheck(faultNode, 180); !b {
			return fmt.Errorf("fault check error,waiting remain %d", remain)
		}
	}

	err := client.Ping(general.RPCAddress(node.NodeAddress))
	if err != nil {
		logE("add peer", "error", err)
		l.nodes.Fault(node)
		return err
	}

	ipfsTimeout, cancelFunc := context.WithTimeout(ctx, time.Duration(l.cfg.Interval)*time.Second)
	var ipfsErr error
	for _, addr := range node.DataStore.Addresses {
		ipfsErr = l.ipfs.SwarmConnect(ipfsTimeout, addr)
		if ipfsErr == nil {
			break
		}
	}
	cancelFunc()
	if ipfsErr != nil {
		logE("add peer", "tag", ipfsErr)
		l.nodes.Fault(node)
		return err
	}
	ethTimeout, cancelFunc := context.WithTimeout(ctx, time.Duration(l.cfg.Interval)*time.Second)
	//fmt.Println("connect eth:", node.Contract.Enode)
	err = l.eth.AddPeer(ethTimeout, node.Contract.Enode)
	if err != nil {
		logE("add peer", "error", err)
		l.nodes.Fault(node)
		return err
	}
	cancelFunc()

	l.nodes.Add(node)
	*result = true
	return nil
}

// Add ...
func (l *BustLinker) Add(r *http.Request, req *core.AddReq, resp *core.AddResp) error {
	if req.AddType == core.AddTypePeer {
		return l.addPeer(r.Context(), &req.Node, &resp.IsSuccess)
	}
	return nil
}

// AddPeer ...
func (l *BustLinker) AddPeer(r *http.Request, info *core.Node, result *bool) error {
	return l.addPeer(r.Context(), info, result)
}

// Peers ...
func (l *BustLinker) Peers(r *http.Request, _ *core.Node, result *[]*core.Node) error {
	l.nodes.Range(func(node *core.Node) bool {
		*result = append(*result, node)
		return true
	})
	return nil
}

func (l *BustLinker) pins(ctx context.Context, result *[]string) error {
	pins, e := l.ipfs.PinLS(ctx)
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
	defer cancelFunc()
	wg.Add(1)
	go func() {
		var err error
		defer func() {
			wg.Done()
			if err != nil && resultErr != nil {
				resultErr <- err
			}
		}()
		if v.PosterHash == "" {
			return
		}

		err = l.nodeConnect(ctx, v.PosterHash)
		if err != nil {
			cancelFunc()
			return
		}
		err = l.ipfs.PinAdd(ctx, v.PosterHash)
		if err != nil {
			cancelFunc()
			return
		}
	}()

	wg.Add(1)
	go func() {
		var err error
		defer func() {
			wg.Done()
			if err != nil && resultErr != nil {
				resultErr <- err
			}
		}()
		if v.ThumbHash == "" {
			return
		}
		err = l.nodeConnect(ctx, v.ThumbHash)
		if err != nil {
			cancelFunc()
			return
		}
		err = l.ipfs.PinAdd(ctx, v.ThumbHash)
		if err != nil {
			cancelFunc()
		}

	}()

	wg.Add(1)
	go func() {
		var err error
		defer func() {
			wg.Done()
			if err != nil && resultErr != nil {
				resultErr <- err
			}
		}()
		if v.SourceHash == "" {
			return
		}
		err = l.nodeConnect(ctx, v.SourceHash)
		if err != nil {
			cancelFunc()
			return
		}
		err = l.ipfs.PinAdd(ctx, v.SourceHash)
		if err != nil {
			cancelFunc()
			return
		}
	}()

	wg.Add(1)
	go func() {
		var err error
		defer func() {
			wg.Done()
			if err != nil && resultErr != nil {
				resultErr <- err
			}
		}()

		if v.M3U8Hash == "" {
			return
		}

		err = l.nodeConnect(ctx, v.M3U8Hash)
		if err != nil {
			cancelFunc()
			return
		}
		err = l.ipfs.PinAdd(ctx, v.M3U8Hash)
		if err != nil {
			cancelFunc()
			return
		}
	}()

	wg.Wait()
	select {
	case e := <-resultErr:
		resultErr = nil
		return e
	default:
	}
	*result = true
	return nil
}

func (l *BustLinker) tagInfo(tag string, info *string) error {
	dTag, e := l.eth.DTag()
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
			resultErr = l.ipfs.SwarmConnect(ctx, addr)
			if resultErr == nil {
				break
			}
		}
	}
	return nil
}
