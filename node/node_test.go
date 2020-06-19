package node

import (
	"fmt"
	"github.com/glvd/accipfs/core"
	"github.com/portmapping/go-reuse"
	"net"
	"testing"
)

type dummyAPI struct {
}

func (d dummyAPI) Ping(req *core.PingReq) (*core.PingResp, error) {
	return nil, nil
}

func (d dummyAPI) ID(req *core.IDReq) (*core.IDResp, error) {
	return &core.IDResp{
		Name:      "abc",
		DataStore: nil,
		Contract:  nil,
	}, nil
}

func TestAcceptNode(t *testing.T) {
	local := &net.TCPAddr{
		IP:   net.IPv4zero,
		Port: 16004,
	}
	listener, err := reuse.ListenTCP("tcp", local)
	if err != nil {
		t.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		node, err := AcceptNode(conn, &dummyAPI{})
		if err != nil {
			continue
		}
		fmt.Println(node.ID())
		//no callback closed
	}

}

func TestConnectToNode(t *testing.T) {
	toNode, err := ConnectToNode(core.Addr{
		Protocol: "tcp",
		IP:       net.IPv4zero,
		Port:     16004,
	}, 0, &dummyAPI{})
	if err != nil {
		return
	}
	toNode.ID()
}
