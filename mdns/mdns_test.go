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

func TestServer_Lookup(t *testing.T) {
	mdns, err := New(config.Default(), func(cfg *OptionConfig) {
		cfg.Service = "_foobar._tcp"
	})
	if err != nil {
		t.Fatal(err)
	}

	s, err := mdns.Server()
	if err != nil {
		t.Fatal(err)
	}
	s.Start()
	defer s.Stop()
	c, err := mdns.Client()
	if err != nil {
		t.Fatal(err)
	}
	entries := make(chan *ServiceEntry, 1)
	var found int32 = 0
	go func() {
		select {
		case e := <-entries:
			if e.Name != "hostname._foobar._tcp.local." {
				t.Fatalf("bad: %v", e)
			}
			if e.Port != 80 {
				t.Fatalf("bad: %v", e)
			}
			if e.Info != "Local web server" {
				t.Fatalf("bad: %v", e)
			}
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
	err = c.Query(params)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if atomic.LoadInt32(&found) == 0 {
		t.Fatalf("record not found")
	}
}
