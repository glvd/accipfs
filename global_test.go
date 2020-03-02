package accipfs_test

import (
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/service"
	"testing"
)

func TestRun(t *testing.T) {
	config.Initialize()
	s, e := service.New(config.Global())
	if e != nil {
		panic(e)
	}
	s.Run()
}
