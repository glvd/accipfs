package main

import (
	"fmt"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/glvd/accipfs/general"
	"github.com/glvd/accipfs/service"
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
			url := fmt.Sprintf("http://localhost:%d/rpc", cfg.Port)
			id, err := service.ID(url)
			if err != nil {
				fmt.Println("local id error:", err.Error())
				return
			}
			for _, addr := range args {
				fmt.Println("connect:", addr)
				remoteURL := fmt.Sprintf("http://%s/rpc", addr)
				remoteNodeInfo := new(core.NodeInfo)
				if err := general.RPCPost(remoteURL, "Accelerate.Connect", id, remoteNodeInfo); err != nil {
					fmt.Println("connect error:", err.Error())
					return
				}

				remoteNodeInfo.RemoteAddr, remoteNodeInfo.Port = general.SplitIP(addr)
				if err := service.AddPeer(url, remoteNodeInfo); err != nil {
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
			if err := general.RPCPost(url, "Accelerate.Peers", &service.Empty{}, reply); err != nil {
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
