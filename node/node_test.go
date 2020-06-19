package node

import (
	"fmt"
	"github.com/glvd/accipfs/core"
	"github.com/portmapping/go-reuse"
	"net"
	"sync"
	"testing"
)

type dummyAPI struct {
	id string
}

func (d dummyAPI) Ping(req *core.PingReq) (*core.PingResp, error) {
	return nil, nil
}

func (d dummyAPI) ID(req *core.IDReq) (*core.IDResp, error) {
	return &core.IDResp{
		Name:      d.id,
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
		fmt.Println("new connector")
		node, err := AcceptNode(conn, &dummyAPI{
			id: "server",
		})
		if err != nil {
			continue
		}
		fmt.Println(node.ID())
		//no callback closed
	}

}

func TestConnectNode(t *testing.T) {
	toNode, err := ConnectNode(core.Addr{
		Protocol: "tcp",
		IP:       net.IPv4zero,
		Port:     16004,
	}, 0, &dummyAPI{
		id: "abc",
	})
	if err != nil {
		return
	}
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			fmt.Println("get id", toNode.ID())
			wg.Done()
		}()
	}
	wg.Wait()
}
