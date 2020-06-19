package main

import (
	"fmt"
	"github.com/glvd/accipfs/basis"
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
			url := config.RPCAddr()

			for _, addr := range args {
				fmt.Printf("connect to [%s]\n", addr)
				req := &core.ConnectToReq{Addr: addr}
				remote, err := client.ConnectTo(url, req)
				if err != nil {
					fmt.Printf("connect error: %v\n", err)
					return
				}

				if err := client.AddPeer(url, remote.Node); err != nil {
					fmt.Printf("add peer error: %v\n", err)
					return
				}
			}

			return
		},
	}
	return connect
}

func nodePeerCmd() *cobra.Command {
	peers := &cobra.Command{
		Use:   "peers",
		Short: "peers run",
		Long:  "show the local node peers",
		Run: func(cmd *cobra.Command, args []string) {
			config.Initialize()
			cfg := config.Global()
			url := fmt.Sprintf("http://localhost:%d/rpc", cfg.API.Port)
			reply := new([]*core.Node)
			node := new(core.Node)
			if err := basis.RPCPost(url, "BustLinker.Peers", node, reply); err != nil {
				fmt.Println("peers error:", err.Error())
			}
			for _, info := range *reply {
				fmt.Println("Peer:", info)
			}
			return
		},
	}
	return peers
}
