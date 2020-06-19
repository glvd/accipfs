package node

import (
	"github.com/glvd/accipfs/config"
	alog "github.com/glvd/accipfs/log"
	"time"
)

var testConfig = config.Default()

func init() {
	alog.InitLog()
	testConfig.Path = ""
	testConfig.Node.BackupSeconds = 3 * time.Second
}
