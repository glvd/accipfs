package main

import (
	"fmt"
	"github.com/glvd/accipfs/client"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/spf13/cobra"
)

func nodeCmd() *cobra.Command {
	nodeCmd := &cobra.Command{
		Use:   "node",
		Short: "node run",
		Long:  "node can operate to change the parameters of some nodes",
	}

	nodeCmd.AddCommand(nodeConnectCmd(), nodePeerCmd())
	return nodeCmd
}

func nodeConnectCmd() *cobra.Command {
	connect := &cobra.Command{
		Use:   "connect",
		Short: "connect run",
		Long:  "connect a remote node",
		Run: func(cmd *cobra.Command, args []string) {
			config.Initialize()
			cfg := config.Global()
			client.InitGlobalClient(&cfg)
			fmt.Printf("connect to [%v]\n", args)
			req := &core.NodeLinkReq{Addrs: args}
			resp, err := client.NodeLink(req)
			if err != nil {
				fmt.Printf("connect error: %v\n", err)
				return
			}
			fmt.Println("success:")
			fmt.Printf("%+v\n", resp.JSON())

			return
		},
	}
	return connect
}

func nodePeerCmd() *cobra.Command {
	peers := &cobra.Command{
		Use:   "list",
		Short: "list run",
		Long:  "show the local node list",
		Run: func(cmd *cobra.Command, args []string) {
			config.Initialize()
			cfg := config.Global()
			client.InitGlobalClient(&cfg)

			list, err := client.List(&core.NodeListReq{})
			if err != nil {
				panic(err)
			}
			for id, info := range list.Nodes {
				fmt.Printf("node:%v\n", id)
				fmt.Printf("info:%v\n", info.JSON())
			}
			return
		},
	}
	return peers
}
