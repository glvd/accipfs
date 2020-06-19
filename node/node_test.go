package node

import (
	"github.com/glvd/accipfs/core"
	"net"
	"testing"
)

func TestAcceptNode(t *testing.T) {

}

func TestConnectToNode(t *testing.T) {
	toNode, err := ConnectToNode(core.Addr{
		Protocol: "tcp",
		IP:       net.IPv4zero,
		Port:     16004,
	}, 16004)
	if err != nil {
		return
	}
	toNode.ID()
}
