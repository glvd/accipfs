package main

import (
	"github.com/glvd/accipfs"
	"github.com/glvd/accipfs/account"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/general"
	"github.com/spf13/cobra"
	"path/filepath"
)

func initCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "init",
		Short: "init run",
		Long:  "init will create the config file with a default settings",
		Run: func(cmd *cobra.Command, args []string) {
			path, err := filepath.Abs(accipfs.DefaultPath)
			if err != nil {
				path = general.CurrentDir()
			}
			config.WorkDir = path
			cfg := config.Default()
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
