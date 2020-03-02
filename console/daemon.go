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
			service.New(config.Global())
		},
	}
}
