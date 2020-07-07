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

			err := config.SaveGenesis(cfg)
			if err != nil {
				panic(err)
			}

			if err := cfg.Init(); err != nil {
				panic(err)
			}
			context := controller.NewContext(cfg)
			c := controller.New(cfg, context)
			if err := c.Initialize(); err != nil {
				panic(err)
			}

			serverConfig, err := config.LoadIPFSServerConfig(cfg)
			if err != nil {
				panic(err)
			}

			acc, err := account.NewAccount(cfg)
			if err != nil {
				panic(err)
			}
			acc.Identity.PeerID = serverConfig.Identity.PeerID
			acc.Identity.PrivKey = serverConfig.Identity.PrivKey

			err = acc.Save(cfg)
			if err != nil {
				panic(err)
			}

		},
	}
	cmd.Flags().StringVar(&restore, "restore", "", "init from a account file")
	return cmd
}
