package main

import (
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/service"
	"github.com/spf13/cobra"
)

func daemonRun() *cobra.Command {
	return &cobra.Command{
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
