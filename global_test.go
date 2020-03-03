package accipfs_test

import (
	"github.com/glvd/accipfs"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/service"
	"testing"
)

func TestRun(t *testing.T) {
	config.Initialize()
	accipfs.DefaultPath = "/mnt/d/ipfstest/ipfsdata"
	s, e := service.New(config.Global())
	if e != nil {
		panic(e)
	}
	s.Run()
}
