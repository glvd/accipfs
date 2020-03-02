package main

import (
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/service"
	"github.com/spf13/cobra"
)

func daemonRun() *cobra.Command {
	return &cobra.Command{
		Use:   "daemon",
		Short: "Run the service as daemon",
		Long:  "Run all the service with a daemon command",
		Run: func(cmd *cobra.Command, args []string) {
			config.Initialize()
			s, e := service.New(config.Global())
			if e != nil {
				panic(e)
			}
			s.Run()
		},
	}
}
