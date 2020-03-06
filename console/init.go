package main

import (
	"github.com/glvd/accipfs"
	"github.com/glvd/accipfs/config"
	"github.com/spf13/cobra"
)

func initCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "init",
		Short: "init run",
		Long:  "init will create the config file with a default settings",
		Run: func(cmd *cobra.Command, args []string) {
			config.SaveConfig(config.Default())
		},
	}
	cmd.Flags().StringVar(&accipfs.DefaultPath, "path", ".", "set work path")
	return cmd
}
