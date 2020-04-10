package main

import (
	"github.com/glvd/accipfs/account"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/controller"
	"github.com/spf13/cobra"
)

func initCmd() *cobra.Command {
	var restore string
	cmd := &cobra.Command{
		Use:   "init",
		Short: "init run",
		Long:  "init will create the config file with a default settings",
		PreRun: func(cmd *cobra.Command, args []string) {

		},
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.Default()
			if err := cfg.Init(); err != nil {
				panic(err)
			}
			c := controller.New(cfg)
			if err := c.Init(); err != nil {
				panic(err)
			}
			acc, err := account.NewAccount(cfg)
			if err != nil {
				panic(err)
			}
			err = acc.Save(cfg)
			if err != nil {
				panic(err)
			}

		},
	}
	cmd.Flags().StringVar(&restore, "restore", "", "init from a account file")
	return cmd
}
