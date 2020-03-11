package service

import (
	"github.com/glvd/accipfs/config"
	"testing"
)

func TestNodeServerETH(t *testing.T) {
	cfg := config.Default()
	node := NewNodeServerETH(*cfg)
	t.Logf("%+v", node)
}
