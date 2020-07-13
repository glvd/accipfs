package node

import (
	"fmt"
	"github.com/glvd/accipfs/basis"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/controller"
	"github.com/glvd/accipfs/core"
	alog "github.com/glvd/accipfs/log"
	ma "github.com/multiformats/go-multiaddr"
	"sync"
	"testing"
	"time"
)

var testConfig = config.Default()

func init() {
	alog.InitLog()
	testConfig.Path = ""
	testConfig.Node.BackupSeconds = 3 * time.Second
}

func TestManager_Store(t *testing.T) {
	cfg := config.Default()
	//context := controller.NewContext(cfg)
	//c := controller.New(cfg)
	//c.API()
	nodeManager := InitManager(cfg)
	multiaddr, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/12345")
	if err != nil {
		panic(err)
	}
	for i := 0; i < 100; i++ {
		connectNode, err := ConnectNode(multiaddr, 0, core.NodeInfo{}, &dummyAPI{})
		if err != nil {
			continue
		}
		nodeManager.Push(connectNode)
		err = nodeManager.Store()
		if err != nil {
			continue
		}
	}
}

func store(wg *sync.WaitGroup, nodeManager core.NodeManager) {
	//defer wg.Done()
	multiaddr, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/12345")
	if err != nil {
		panic(err)
	}
	for i := 0; i < 1000; i++ {
		connectNode, err := ConnectNode(multiaddr, 0, core.NodeInfo{}, &dummyAPI{
			id: basis.UUID(),
		})

		if err != nil {
			fmt.Println("error", err)
		}
		fmt.Println("remote id:", connectNode.ID())
		nodeManager.Push(connectNode)
		err = nodeManager.Store()
		if err != nil {
			fmt.Println("error", err)
		}
	}
}

func load(wg *sync.WaitGroup, nodeManager core.NodeManager) {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		err := nodeManager.Load()
		if err != nil {
			fmt.Println("error", err)
			continue
		}
		nodeManager.Range(func(key string, n core.Node) bool {
			fmt.Println("key:", key, "node", n.ID())
			return true
		})
	}
}

func TestManager_Load(t *testing.T) {
	cfg := config.Default()
	controller.New(cfg)
	InitManager(cfg)
	//nodeManager := Manager(cfg, ctx)
	defer GlobalManager.Close()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go load(wg, GlobalManager)
	store(wg, GlobalManager)
	wg.Done()
}
