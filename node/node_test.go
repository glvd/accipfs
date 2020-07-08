package node

import (
	"fmt"
	"github.com/glvd/accipfs/basis"
	"github.com/glvd/accipfs/core"
	"github.com/godcong/scdt"
	ma "github.com/multiformats/go-multiaddr"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
	"testing"
)

type dummyAPI struct {
	id string
}

func (d *dummyAPI) Add(req *core.AddReq) (*core.AddResp, error) {
	panic("implement me")
}

func (d *dummyAPI) NodeAPI() core.NodeAPI {
	panic("implement me")
}

func (d *dummyAPI) AddrInfo(req *core.AddrReq) (*core.AddrResp, error) {
	return nil, nil
}

func init() {
	runtime.GOMAXPROCS(2)
}

func (d dummyAPI) Ping(req *core.PingReq) (*core.PingResp, error) {
	return nil, nil
}

func (d dummyAPI) ID(req *core.IDReq) (*core.IDResp, error) {
	return &core.IDResp{
		ID: d.id,
	}, nil
}

func TestAcceptNode(t *testing.T) {
	listener, err := scdt.NewListener("0.0.0.0:12345")
	if err != nil {
		t.Fatal(err)
	}
	listener.ID(func() string {
		return basis.UUID()
	})
	listener.HandleRecv(func(id string, message *scdt.Message) ([]byte, bool) {
		fmt.Println("id:", id, "message:", message, "data:", string(message.Data))
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
	multiaddr, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/12345")
	if err != nil {
		return
	}
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			toNode, err := ConnectNode(multiaddr, 0, &dummyAPI{
				id: fmt.Sprintf("id(%v),client request", 0),
			})
			if err != nil {
				t.Fatal(err)
			}
			j := 0
			for ; j < 10; j++ {
				toNode.ID()
			}
			for ; j < 10; j++ {
				toNode.Info()

			}

			info, err := toNode.Info()

			fmt.Println("get id", i, "index", j, "id", toNode.ID(), "info", info.JSON())

			err = toNode.Close()
			if err != nil {
				fmt.Println("err", err)
			}
		}(i)
	}
	wg.Wait()
}
