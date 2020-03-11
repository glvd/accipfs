package accipfs_test

import (
	"github.com/glvd/accipfs"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/service"
	"testing"
)

func TestRun(t *testing.T) {
	config.Initialize()
	accipfs.DefaultPath = "data"
	s, e := service.NewClient(config.Global())
	if e != nil {
		panic(e)
	}
	s.Run()
}
