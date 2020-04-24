package main

import (
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/mdns"
	"github.com/glvd/accipfs/service"
	"github.com/spf13/cobra"
)

func daemonCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "daemon",
		Short: "Run the service as daemon",
		Long:  "Run all the service with a daemon command",
		Run: func(cmd *cobra.Command, args []string) {
			config.Initialize()
			cfg := config.Global()
			s, e := service.New(&cfg)
			if e != nil {
				panic(e)
			}
			dns, e := mdns.New(&cfg, func(c *mdns.OptionConfig) {
				c.RegisterLocalIP(&cfg)
				c.Service = "_bustlinker._udp"
			})
			if e != nil {
				panic(e)
			}
			server, e := dns.Server()
			if e != nil {
				panic(e)
			}
			defer func() {
				server.Stop()
			}()
			server.Start()
			defer func() {
				if err := s.Stop(); err != nil {
					panic(err)
				}
			}()
			if err := s.Start(); err != nil {
				panic(err)
			}
		},
	}
}
