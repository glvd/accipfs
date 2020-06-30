package service

import (
	"fmt"
	"github.com/glvd/accipfs/account"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/controller"
	"github.com/glvd/accipfs/core"
	"github.com/glvd/accipfs/node"
	"github.com/glvd/accipfs/task"
	"sync"

	"github.com/robfig/cron/v3"
	"go.uber.org/atomic"
)

// BustLinker ...
type BustLinker struct {
	id         core.Node
	nodes      core.NodeManager
	tasks      task.Task
	lock       *atomic.Bool
	self       *account.Account
	cfg        *config.Config
	cron       *cron.Cron
	listener   core.Listener
	controller *controller.Controller
}

// NewBustLinker ...
func NewBustLinker(cfg *config.Config) (linker *BustLinker, err error) {
	linker = &BustLinker{
		lock: atomic.NewBool(false),
		cfg:  cfg,
	}

	selfAcc, err := account.LoadAccount(cfg)
	if err != nil {
		return nil, err
	}
	linker.self = selfAcc
	linker.controller = controller.New(cfg)
	linker.nodes = node.New(cfg, linker.controller)
	linker.listener = NewLinkListener(cfg, linker.nodes.HandleConn)

	return linker, nil
}

// Start ...
func (l *BustLinker) Start() {
	l.controller.Run()
	go l.listener.Listen()
}

// Run ...
func (l *BustLinker) Run() {
	if l.lock.Load() {
		fmt.Println(module, "the previous task has not been completed")
		return
	}
	l.lock.Store(true)
	defer l.lock.Store(false)
	wg := &sync.WaitGroup{}
	l.nodes.Range(func(key string, node core.Node) bool {
		return true
	})
	wg.Wait()
	fmt.Println(module, "syncing done")
}

// WaitingForReady ...
func (l *BustLinker) WaitingForReady() {

}

// Stop ...
func (l *BustLinker) Stop() {
	ctx := l.cron.Stop()
	<-ctx.Done()
}
