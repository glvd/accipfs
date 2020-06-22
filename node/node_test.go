package node

import (
	"fmt"
	"github.com/glvd/accipfs/core"
	"github.com/panjf2000/ants/v2"
	"github.com/portmapping/go-reuse"
	"net"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
	"testing"
	"time"
)

type dummyAPI struct {
	id string
}

func init() {
	runtime.GOMAXPROCS(2)
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
	go func() {
		ip := "0.0.0.0:6060"
		if err := http.ListenAndServe(ip, nil); err != nil {
			fmt.Printf("start pprof failed on %s\n", ip)
		}
	}()
	pool, _ := ants.NewPool(ants.DefaultAntsPoolSize, ants.WithNonblocking(false))

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		pool.Submit(func() {
			node, err := AcceptNode(conn, &dummyAPI{
				id: "server",
			})
			if err != nil {
				fmt.Println("err", err)
				return
			}
			fmt.Println(node.ID())
			node.Closed(func(n core.Node) {
				node = nil
			})
		})
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
				wg.Done()
				t.Fatal(err)
			}
			j := 0
			for ; j < 10; j++ {
				toNode.ID()
			}
			fmt.Println("get id", i, "index", j, toNode.ID())
			time.Sleep(30 * time.Minute)
			err = toNode.Close()
			if err != nil {
				fmt.Println("err", err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
