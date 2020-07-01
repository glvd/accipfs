package node

import (
	"fmt"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/controller"
	"github.com/glvd/accipfs/core"
	alog "github.com/glvd/accipfs/log"
	ma "github.com/multiformats/go-multiaddr"
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
	c := controller.New(cfg)
	nodeManager := Manager(cfg, c.GetAPI())
	multiaddr, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/12345")
	if err != nil {
		panic(err)
	}
	for i := 0; i < 100; i++ {
		connectNode, err := ConnectNode(multiaddr, 0, &dummyAPI{})
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
func TestManager_Load(t *testing.T) {
	cfg := config.Default()
	c := controller.New(cfg)
	nodeManager := Manager(cfg, c.GetAPI())
	for i := 0; i < 100; i++ {
		err := nodeManager.Load()
		if err != nil {
			continue
		}
		nodeManager.Range(func(key string, n core.Node) bool {
			fmt.Println("key:", key, "node", n.ID())
			return true
		})
	}
}
