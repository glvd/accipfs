package node

import (
	"github.com/glvd/accipfs/config"
	alog "github.com/glvd/accipfs/log"
	"testing"
	"time"
)

var testConfig = config.Default()

func init() {
	alog.InitLog()
	testConfig.Path = ""
	testConfig.Node.BackupSeconds = 3 * time.Second
}

func TestManager_Load(t *testing.T) {
	cfg := config.Default()
	nodeManager := New(cfg, nil)
	err := nodeManager.Load()
	if err != nil {
		panic(err)
	}
}
