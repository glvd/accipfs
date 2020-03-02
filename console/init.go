package main

import (
	"github.com/glvd/accipfs/config"
	"github.com/spf13/cobra"
)

func initCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "init run",
		Long:  "init will create the config file with a default settings",
		Run: func(cmd *cobra.Command, args []string) {
			config.SaveConfig(config.Default())
		},
	}
}
