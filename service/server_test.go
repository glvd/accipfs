package service

import (
	"github.com/glvd/accipfs/config"
	"testing"
)

func TestNodeServerETH(t *testing.T) {
	config.WorkDir = "D:\\workspace\\pvt"
	c := config.Default()
	eth := NewNodeServerETH(*c)
	t.Logf("%+v,cfg:%+v", eth, *c)
	if err := eth.Init(); err != nil {
		t.Error(err)
		return
	}

}
