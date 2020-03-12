package service

import (
	"github.com/glvd/accipfs/config"
	"testing"
)

func TestNodeServerETH(t *testing.T) {
	c := config.Default()
	c.Path = "D:\\workspace\\pvt"
	eth := NewNodeServerETH(*c)
	t.Logf("%+v", eth)
}
