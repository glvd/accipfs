package main

import (
	"fmt"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/log"
	"github.com/glvd/accipfs/mdns"
	"github.com/glvd/accipfs/service"
	"github.com/spf13/cobra"
	"time"
)

func daemonCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "daemon",
		Short: "Run the service as daemon",
		Long:  "Run all the service with a daemon command",
		Run: func(cmd *cobra.Command, args []string) {
			log.InitLog()

			config.Initialize()
			cfg := config.Global()
			linker, err := service.NewBustLinker(&cfg)
			if err != nil {
				panic(err)
			}
			dns, e := mdns.New(&cfg, func(c *mdns.OptionConfig) {
				c.Service = "_bustlinker._udp"
				c.RegisterLocalIP(&cfg)
			})
			if e != nil {
				panic(e)
			}
			server, e := dns.Server()
			if e != nil {
				fmt.Printf("multicast group join failed(%v)", e)
			}
			defer func() {
				server.Stop()
			}()
			server.Start()
			defer func() {
				linker.Stop()
			}()
			linker.Start()
			time.Sleep(30 * time.Minute)
		},
	}
}
