package main

import (
	"fmt"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/log"
	"github.com/glvd/accipfs/mdns"
	"github.com/spf13/cobra"
	"net"
	"os"
	"os/signal"
	"time"
)

var rootCmd = &cobra.Command{}

func main() {
	log.InitLog()
	var client bool
	var service string
	port := 8080
	info := "accipfs local server"
	rootCmd.PersistentFlags().BoolVar(&client, "client", false, "enable client model")
	rootCmd.PersistentFlags().StringVar(&service, "service", "_foobar._tcp", "set the service name")

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		fmt.Println("mdns test running")
		m, err := mdns.New(config.Default(), func(cfg *mdns.OptionConfig) {
			cfg.Service = service
			addrs, err := net.InterfaceAddrs()
			if err != nil {
				return
			}
			for i := range addrs {
				if ipnet, ok := addrs[i].(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						cidr, _, err := net.ParseCIDR(addrs[i].String())
						if err == nil {
							fmt.Println("register ip addr:", cidr.String())
							cfg.IPs = append(cfg.IPs, cidr)
						}
					}
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
			defer s2.Stop()
			handler := make(chan os.Signal, 1)
			signal.Notify(handler, os.Interrupt)
			for sig := range handler {
				if sig == os.Interrupt {
					time.Sleep(1e9)
					break
				}
			}
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
			//var found int32
			go func() {
				for e := range entries {
					log.Module("main").Infow("output detail", "name", e.Name, "host", e.Host, "fields", e.InfoFields, "ipv4", e.AddrV4, "ipv6", e.AddrV6)
					log.Module("main").Infow("output addr", "addr", e.Addr.String())
					log.Module("main").Infow("output port", "port", e.Port, "want", port)
					log.Module("main").Infow("output info", "info", e.Info, "want", info)
				}
				//select {
				//case e := <-entries:
				//	//if e.Name != "accipfs._foobar._tcp.local." {
				//	//	log.Module("main").Fatalf("bad: %v", e)
				//	//}
				//	log.Module("main").Infow("output detail", "name", e.Name, "host", e.Host, "fields", e.InfoFields, "ipv4", e.AddrV4.String(), "ipv6", e.AddrV6.String())
				//	log.Module("main").Infow("output addr", "addr", e.Addr.String())
				//	log.Module("main").Infow("output port", "port", e.Port, "want", port)
				//	log.Module("main").Infow("output info", "info", e.Info, "want", info)
				//
				//	atomic.StoreInt32(&found, 1)
				//
				//case <-time.After(80 * time.Second):
				//	log.Module("main").Fatalf("timeout")
				//}
			}()

			params := &mdns.QueryParam{
				Service: service,
				Domain:  "local",
				Timeout: 3 * time.Second,
				Entries: entries,
			}
			err = c.Query(params)
			if err != nil {
				log.Module("main").Fatalf("err: %v\n", err)
			}
			//if atomic.LoadInt32(&found) == 0 {
			//	log.Module("main").Fatalf("record not found")
			//}
			time.Sleep(5 * time.Second)
		}
	}
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
