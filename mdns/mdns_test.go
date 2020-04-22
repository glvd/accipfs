package mdns

import (
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/log"
	"net"
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
	mdns, err := New(config.Default(), func(cfg *OptionConfig) {
		cfg.Service = "_foobar._tcp"
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			t.Log(err)
			return
		}
		for i := range addrs {
			if ipnet, ok := addrs[i].(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					cidr, _, err := net.ParseCIDR(addrs[i].String())
					if err == nil {
						cfg.IPs = append(cfg.IPs, cidr)
					}
					cfg.IPs = append(cfg.IPs, cidr)
				}
			}
		}
		//cfg.IPs = append(cfg.IPs, net.ParseIP("192.168.1.45"), net.ParseIP("192.168.1.13"))
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
			log.Module("main").Infow("output detail", "name", e.Name, "host", e.Host, "fields", e.InfoFields, "ipv4", e.AddrV4.String(), "ipv6", e.AddrV6.String())
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
	err = c.Query(params)
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
