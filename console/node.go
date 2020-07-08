package main

import (
	"fmt"
	"github.com/glvd/accipfs/basis"
	"github.com/glvd/accipfs/client"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	ma "github.com/multiformats/go-multiaddr"
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
			var addrs []ma.Multiaddr
			for _, addr := range args {
				fmt.Printf("connect to [%s]\n", addr)
				multiaddr, err := ma.NewMultiaddr(addr)
				if err != nil {
					fmt.Printf("wrong connect address error: %v\n", err)
					return
				}
				addrs = append(addrs, multiaddr)
				req := &core.NodeLinkReq{Addrs: addrs}
				resp, err := client.NodeLink(req)
				if err != nil {
					fmt.Printf("connect error: %v\n", err)
					return
				}
				fmt.Println("success:")
				fmt.Printf("%+v", resp)
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
