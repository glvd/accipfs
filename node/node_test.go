package node

import (
	"fmt"
	"github.com/godcong/scdt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/glvd/accipfs/core"
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
	listener, err := scdt.NewListener("0.0.0.0:12345")
	if err != nil {
		t.Fatal(err)
	}
	listener.HandleRecv(func(id string, message *scdt.Message) ([]byte, bool) {
		fmt.Println("id:", id, "message", message)
		return []byte("success"), true
	})

	ip := "0.0.0.0:6060"
	if err := http.ListenAndServe(ip, nil); err != nil {
		fmt.Printf("start pprof failed on %s\n", ip)
	}
	listener.Stop()
}

func TestConnectNode(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			toNode, err := ConnectNode(core.Addr{
				Protocol: "tcp",
				IP:       net.IPv4zero,
				Port:     12345,
			}, 0, &dummyAPI{
				id: fmt.Sprintf("id(%v),client request", 0),
			})
			if err != nil {
				wg.Done()
				t.Fatal(err)
			}
			j := 0
			for ; j < 100; j++ {
				toNode.ID()
			}
			for ; j < 100; j++ {
				toNode.Info()
			}
			fmt.Println("get id", i, "index", j, "id", toNode.ID(), "info", toNode.Info())
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
