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
	testConfig.Node.BackupSeconds = 5 * 60 * time.Second
}

func TestNew(t *testing.T) {
	m := New(testConfig)
	err := m.Load()
	if err != nil {
		t.Log(err)
	}
}
