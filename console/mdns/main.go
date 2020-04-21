package main

import (
	"fmt"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/log"
	"github.com/glvd/accipfs/mdns"
	"net"
	"time"
)

func main() {
	log.InitLog()
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
	})
	if err != nil {
		panic(err)
	}

	s2, err := m.Server()
	if err != nil {
		panic(err)
	}
	s2.Start()
	time.Sleep(5 * time.Minute)
	defer s2.Stop()
}
