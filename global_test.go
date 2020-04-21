package accipfs_test

import (
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/service"
	"testing"
)

func TestRun(t *testing.T) {
	config.WorkDir = "D:\\workspace\\pvt"
	config.Initialize()
	cfg := config.Global()
	s, e := service.New(&cfg)
	t.Log("accelerate new")
	if e != nil {
		panic(e)
	}
	defer s.Stop()
	s.Start()
}
