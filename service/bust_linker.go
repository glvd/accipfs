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
	linker.manager = node.Manager(cfg, linker.controller.GetAPI())
	linker.listener = NewLinkListener(cfg, linker.manager.Conn)

	return linker, nil
}

// Start ...
func (l *BustLinker) Start() {
	l.controller.Run()
	go l.listener.Listen()
}

// Run ...
func (l *BustLinker) Run() {
	if !l.lock.CAS(false, true) {
		return
	}
	defer l.lock.Store(false)
	//todo
}

// WaitingForReady ...
func (l *BustLinker) WaitingForReady() {

}

// Stop ...
func (l *BustLinker) Stop() {
}
