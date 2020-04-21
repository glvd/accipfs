package main

import (
	"fmt"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/log"
	"github.com/glvd/accipfs/mdns"
	"github.com/spf13/cobra"
	"net"
	"sync/atomic"
	"time"
)

var rootCmd = &cobra.Command{}

func main() {
	var client bool
	port := 8080
	info := "accipfs local server"
	rootCmd.PersistentFlags().BoolVar(&client, "client", false, "enable client model")

	log.InitLog()

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		fmt.Println("mdns test running")
		m, err := mdns.New(config.Default(), func(cfg *mdns.OptionConfig) {
			cfg.Service = "_foobar._tcp"
			addrs, err := net.InterfaceAddrs()
			if err != nil {
				return
			}
			for i := range addrs {
				cidr, _, err := net.ParseCIDR(addrs[i].String())
				if err == nil {
					fmt.Println("ip added:", addrs[i].String())
					cfg.IPs = append(cfg.IPs, cidr)
				}
			}
			cfg.Port = uint16(port)
		})
		if err != nil {
			panic(err)
		}

		if !client {
			s2, err := m.Server()
			if err != nil {
				panic(err)
			}
			s2.Start()
			time.Sleep(5 * time.Minute)
			defer s2.Stop()
		} else {
			c, err := m.Client()
			if err != nil {
				panic(err)
			}
			entries := make(chan *mdns.ServiceEntry, 1)

			//err = c.Lookup("_foobar._tcp", entries)
			//if err != nil {
			//	t.Log(err)
			//}
			//fmt.Printf("entries:%+v", entries)

			var found int32
			go func() {
				select {
				case e := <-entries:
					if e.Name != "hostname._foobar._tcp.local." {
						log.Module("main").Fatalf("bad: %v", e)
					}
					log.Module("main").Infow("name", e.Name, "host", e.Host, "fields", e.InfoFields, "ipv4", e.AddrV4.String(), "ipv6", e.AddrV6.String())
					log.Module("main").Infow(e.Addr.String())
					log.Module("main").Infow("port", e.Port, "want", port)
					log.Module("main").Infow("info", e.Info, "want", info)

					atomic.StoreInt32(&found, 1)

				case <-time.After(80 * time.Second):
					log.Module("main").Fatalf("timeout")
				}
			}()

			params := &mdns.QueryParam{
				Service: "_foobar._tcp",
				Domain:  "local",
				Timeout: 50 * time.Millisecond,
				Entries: entries,
			}
			err = c.Query(params)
			if err != nil {
				log.Module("main").Fatalf("err: %v\n", err)
			}
			if atomic.LoadInt32(&found) == 0 {
				log.Module("main").Fatalf("record not found")
			}
		}
	}
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
