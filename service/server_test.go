package service

import (
	"github.com/glvd/accipfs/config"
	"testing"
)

func TestNodeServerETH(t *testing.T) {
	config.WorkDir = "D:\\workspace\\pvt"
	err := config.SaveConfig(config.Default())
	if err != nil {
		t.Error(err)
		return
	}
	config.Initialize()
	eth := NewNodeServerETH(config.Global())
	t.Logf("%+v", eth)
	if err := eth.Init(); err != nil {
		t.Error(err)
		return
	}
	ipfs := NewNodeServerIPFS(config.Global())
	t.Logf("%+v", ipfs)
	if err := ipfs.Init(); err != nil {
		return
	}
}
