package main

import (
	"github.com/glvd/accipfs"
	"github.com/glvd/accipfs/account"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/general"
	"github.com/glvd/accipfs/service"
	"github.com/goextension/log/zap"
	"github.com/spf13/cobra"
	"path/filepath"
)

func initCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "init",
		Short: "init run",
		Long:  "init will create the config file with a default settings",
		Run: func(cmd *cobra.Command, args []string) {
			zap.InitZapFileSugar()
			path, err := filepath.Abs(accipfs.DefaultPath)
			if err != nil {
				path = general.CurrentDir()
			}
			config.WorkDir = path
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
	cmd.Flags().StringVar(&accipfs.DefaultPath, "path", ".", "set work path")
	return cmd
}
