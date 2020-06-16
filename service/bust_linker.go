package service

import (
	"github.com/glvd/accipfs/account"
	"github.com/glvd/accipfs/cache"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/controller"
	"github.com/glvd/accipfs/core"
	"github.com/glvd/accipfs/task"
	"net"
	"sync"

	"github.com/robfig/cron/v3"
	"go.uber.org/atomic"
)

// BustLinker ...
type BustLinker struct {
	id         core.Node
	nodes      core.NodeManager
	tasks      task.Task
	hashes     cache.HashCache
	lock       *atomic.Bool
	self       *account.Account
	cfg        *config.Config
	cron       *cron.Cron
	listener   core.Listener
	controller *controller.Controller
	api        *controller.API
}

// NewBustLinker ...
func NewBustLinker(cfg *config.Config) (linker *BustLinker, err error) {
	linker = &BustLinker{
		hashes: cache.NewHashCache(cfg),
		lock:   atomic.NewBool(false),
		cfg:    cfg,
	}

	selfAcc, err := account.LoadAccount(cfg)
	if err != nil {
		return nil, err
	}
	linker.self = selfAcc
	linker.listener = NewLinkListener(cfg, linker.cb)
	linker.controller = controller.New(cfg)
	return linker, nil
}

// Start ...
func (l *BustLinker) Start() {
	go l.listener.Listen()
	l.controller.Run()
}

// Run ...
func (l *BustLinker) Run() {
	if l.lock.Load() {
		output("bust linker", "the previous task has not been completed")
		return
	}
	l.lock.Store(true)
	defer l.lock.Store(false)
	wg := &sync.WaitGroup{}
	l.nodes.Range(func(key string, node core.Node) bool {
		return true
	})
	wg.Wait()
	output("bust linker", "syncing done")
}

// WaitingForReady ...
func (l *BustLinker) WaitingForReady() {

}

// Stop ...
func (l *BustLinker) Stop() {
	ctx := l.cron.Stop()
	<-ctx.Done()
}

func (l *BustLinker) cb(conn net.Conn) {
	//todo:new node
}
