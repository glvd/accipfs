package main

import (
	"fmt"
	"github.com/glvd/accipfs/client"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/glvd/accipfs/general"
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
				fmt.Printf("connect to [%s]\n", url.String())

				remoteNodeInfo := new(core.NodeInfo)
				if err := general.RPCPost(url.String(), "Accelerate.ConnectTo", addr, remoteNodeInfo); err != nil {
					fmt.Println("connect error:", err)
					return
				}

				if err := client.AddPeer(url.String(), remoteNodeInfo); err != nil {
					fmt.Println("add peer error:", err)
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
			url := fmt.Sprintf("http://localhost:%d/rpc", cfg.Port)
			reply := new([]*core.NodeInfo)
			if err := general.RPCPost(url, "Accelerate.Peers", &core.Empty{}, reply); err != nil {
				fmt.Println("peers error:", err.Error())
			}
			for _, info := range *reply {
				fmt.Println("Peer:", info.Name)
			}
			return
		},
	}
	return peers
}
