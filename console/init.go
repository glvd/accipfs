package main

import (
	"github.com/glvd/accipfs/account"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/controller"
	"github.com/glvd/accipfs/node"
	ipfsCfg "github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"github.com/spf13/cobra"
	"os"
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
			context := node.NewContext(cfg)
			c := controller.New(cfg, context)
			if err := c.Initialize(); err != nil {
				panic(err)
			}
			//file, err := os.OpenFile("key", os.O_CREATE|os.O_SYNC|os.O_RDWR|os.O_TRUNC, 0755)
			//if err != nil {
			//	panic(err)
			//}

			identity, err := ipfsCfg.CreateIdentity(os.Stdout, []options.KeyGenerateOption{options.Key.Type("ed25519")})
			if err != nil {
				panic(err)
			}
			cfg.Identity = identity.PeerID
			cfg.PrivateKey = identity.PrivKey
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
