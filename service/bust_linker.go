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
	context := controller.NewContext(cfg)

	linker.controller = controller.New(cfg)

	////todo
	//info, err := context.NodeAddrInfo(&core.AddrReq{})
	//if err != nil {
	//	return nil
	//}
	//m.local = core.NodeInfo{
	//	AddrInfo:        *info.AddrInfo,
	//	AgentVersion:    "",
	//	ProtocolVersion: "",
	//}

	linker.manager = node.InitManager(cfg)
	linker.listener = newLinkListener(cfg, linker.manager.Conn)

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
	//todo do something on run
}

// WaitingForReady ...
func (l *BustLinker) WaitingForReady() {

}

// Stop ...
func (l *BustLinker) Stop() {
}
