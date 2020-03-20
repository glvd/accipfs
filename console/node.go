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
	var addr string
	connect := &cobra.Command{
		Use:   "connect",
		Short: "connect run",
		Long:  "connect a remote node",
		Run: func(cmd *cobra.Command, args []string) {
			config.Initialize()
			cfg := config.Global()
			url := fmt.Sprintf("http://localhost:%d/rpc", cfg.Port)
			reply := new(core.NodeInfo)
			if err := general.RPCPost(url, "Accelerate.ID", &service.Empty{}, reply); err != nil {
				fmt.Println("local id error:", err.Error())
			}
			remoteURL := fmt.Sprintf("http://%s/rpc", addr)
			remoteNodeInfo := new(core.NodeInfo)
			if err := general.RPCPost(remoteURL, "Accelerate.Connect", reply, remoteNodeInfo); err != nil {
				fmt.Println("connect error:", err.Error())
			}

			if err := general.RPCPost(remoteURL, "Accelerate.AddPeer", &service.Empty{}, remoteNodeInfo); err != nil {
				fmt.Println("remote id error:", err.Error())
			}

			return
		},
	}
	connect.Flags().StringVar(&addr, "addr", "localhost:20304", "set a remote address to connect")
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
