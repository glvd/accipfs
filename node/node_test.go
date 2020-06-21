package node

import (
	"fmt"
	"github.com/glvd/accipfs/core"
	"github.com/portmapping/go-reuse"
	"net"
	"runtime"
	"sync"
	"testing"
)

type dummyAPI struct {
	id string
}

func init() {
	runtime.GOMAXPROCS(4)
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
		node, err := AcceptNode(conn, &dummyAPI{
			id: "server",
		})
		if err != nil {
			fmt.Println("err", err)
			continue
		}
		go func() {
			fmt.Println(node.ID())
		}()

		//no callback closed
	}

}

func TestConnectNode(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			toNode, err := ConnectNode(core.Addr{
				Protocol: "tcp",
				IP:       net.IPv4zero,
				Port:     16004,
			}, 0, &dummyAPI{
				id: fmt.Sprintf("id(%v),client request", 0),
			})
			if err != nil {
				t.Fatal(err)
			}
			j := 0
			//for ; j < 100; j++ {
				fmt.Println("get id", i, "index", j, toNode.ID())
			//}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
