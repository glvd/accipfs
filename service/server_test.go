package service

import (
	"github.com/glvd/accipfs/config"
	"github.com/goextension/log/zap"
	"testing"
	"time"
)

func init() {
	zap.InitZapFileSugar()
}
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
	//if err := eth.Init(); err != nil {
	//	t.Error(err)
	//	return
	//}
	if err := eth.Start(); err != nil {
		t.Error(err)
		return
	}

	ipfs := NewNodeServerIPFS(config.Global())
	t.Logf("%+v", ipfs)
	//if err := ipfs.Init(); err != nil {
	//	t.Error(err)
	//	return
	//}
	if err := ipfs.Start(); err != nil {
		t.Error(err)
		return
	}

	time.Sleep(1 * time.Minute)

	if err := eth.Stop(); err != nil {
		t.Error(err)
		return
	}

	if err := ipfs.Stop(); err != nil {
		t.Error(err)
		return
	}
	t.Log("done")
}
