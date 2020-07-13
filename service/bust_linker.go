package service

import (
	"github.com/glvd/accipfs/account"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/controller"
	"github.com/glvd/accipfs/core"
	"github.com/glvd/accipfs/node"
	"github.com/glvd/accipfs/task"
	"go.uber.org/atomic"
)

// BustLinker ...
type BustLinker struct {
	id         core.Node
	manager    core.NodeManager
	tasks      task.Task
	lock       *atomic.Bool
	self       *account.Account
	cfg        *config.Config
	listener   core.Listener
	controller *controller.Controller
	api        *APIContext
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
	linker.manager, err = node.InitManager(cfg)

	linker.api = NewAPIContext(cfg, linker.manager, linker.controller)

	linker.listener = newLinkListener(cfg, linker.manager.Conn)
	return linker, nil
}

// Start ...
func (l *BustLinker) Start() {
	l.controller.Run()
	l.api.Start()
	go l.listener.Listen()
	l.afterStart()
}

// Run ...
func (l *BustLinker) Run() {
	if !l.lock.CAS(false, true) {
		return
	}
	defer l.lock.Store(false)
	//todo do something on run
}

// WaitingForReady ...
func (l *BustLinker) WaitingForReady() {

}

// Stop ...
func (l *BustLinker) Stop() {
}

func (l *BustLinker) afterStart() {
	info, err := l.api.NodeAddrInfo(&core.AddrReq{})
	if err != nil {
		return
	}

	l.manager.Local().Update(func(data *core.LocalData) {
		data.Node.AddrInfo = info.AddrInfo
		//log.Infow("debug", "data", data.Node.DataStore, "info", info.AddrInfo.DataStore)
	})
}
