package main

import (
	"github.com/glvd/accipfs"
	"github.com/glvd/accipfs/account"
	"github.com/glvd/accipfs/config"
	"github.com/spf13/cobra"
)

func initCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "init",
		Short: "init run",
		Long:  "init will create the config file with a default settings",
		Run: func(cmd *cobra.Command, args []string) {
			config.WorkDir = accipfs.DefaultPath
			cfg := config.Default()
			config.Initialize()
			acc, e := account.NewAccount(cfg)
			if e != nil {
				panic(e)
			}
			e = acc.Save(cfg)
			if e != nil {
				panic(e)
			}
			//config.SaveConfig(config.Default())
		},
	}
	cmd.Flags().StringVar(&accipfs.DefaultPath, "path", ".", "set work path")
	return cmd
}
