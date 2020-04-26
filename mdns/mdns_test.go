package mdns

import (
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/log"
	"sync/atomic"
	"testing"
	"time"
)

func init() {
	log.InitLog()
}

func TestMulticastDNS_Server(t *testing.T) {
	mdns, err := New(config.Default())
	if err != nil {
		t.Fatal(err)
	}

	s2, err := mdns.Server()
	if err != nil {
		t.Fatal(err)
	}
	s2.Start()
	defer s2.Stop()
}

func TestMulticastDNS_Lookup(t *testing.T) {
	c := config.Default()
	mdns, err := New(c, func(cfg *OptionConfig) {
		//cfg.Instance = "test"
		//cfg.RegisterLocalIP(c)
	})
	if err != nil {
		t.Fatal(err)
	}

	//s, err := mdns.Server()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//s.Start()
	//defer s.Stop()

	cli, err := mdns.Client()
	if err != nil {
		t.Fatal(err)
	}

	entries := make(chan *ServiceEntry, 1)

	//err = c.Lookup("_foobar._tcp", entries)
	//if err != nil {
	//	t.Log(err)
	//}
	//fmt.Printf("entries:%+v", entries)

	var found int32
	go func() {
		select {
		case e := <-entries:
			t.Log("entries")
			if e.Name != "accipfs._foobar._tcp.local." {
				log.Module("main").Fatalf("bad: %v", e)
			}
			log.Module("main").Infow("output detail", "name", e.Name, "host", e.Host, "fields", e.InfoFields)
			log.Module("main").Infow("output addr", "addr", e.Addr.String())
			log.Module("main").Infow("output port", "port", e.Port, "want", 0)
			log.Module("main").Infow("output info", "info", e.Info, "want", 0)

			atomic.StoreInt32(&found, 1)

		case <-time.After(80 * time.Millisecond):
			t.Fatalf("timeout")
		}
	}()

	params := &QueryParam{
		Service: "_foobar._tcp",
		Domain:  "local",
		Timeout: 50 * time.Millisecond,
		Entries: entries,
	}
	err = cli.Query(params)
	if err != nil {
		t.Fatalf("err: %v\n", err)
	}
	//err = c.Lookup("_foobar._tcp", entries)
	//if err != nil {
	//	t.Fatalf("err: %v\n", err)
	//}
	if atomic.LoadInt32(&found) == 0 {
		t.Fatalf("record not found")
	}
}
func TestFoo(t *testing.T) {
	t.Log(interceptAccountName("0x427a762a04b3c3ac26140125236167bd339f5a5c"))
}
