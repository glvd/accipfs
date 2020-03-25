package main

import (
	"github.com/glvd/accipfs/account"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/service"
	"github.com/spf13/cobra"
)

func initCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "init",
		Short: "init run",
		Long:  "init will create the config file with a default settings",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.Default()
			if err := cfg.Init(); err != nil {
				panic(err)
			}
			ipfs := service.NewNodeServerIPFS(cfg)
			if err := ipfs.Init(); err != nil {
				panic(err)
			}
			eth := service.NewNodeServerETH(cfg)
			if err := eth.Init(); err != nil {
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
			//config.SaveConfig(config.Default())
		},
	}
	return cmd
}
