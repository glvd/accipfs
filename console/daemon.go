package main

import (
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/service"
	"github.com/goextension/log/zap"
	"github.com/spf13/cobra"
)

func daemonCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "daemon",
		Short: "Run the service as daemon",
		Long:  "Run all the service with a daemon command",
		Run: func(cmd *cobra.Command, args []string) {
			zap.InitZapFileSugar()
			config.Initialize()
			cfg := config.Global()
			s, e := service.NewRPCServer(&cfg)
			if e != nil {
				panic(e)
			}
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
