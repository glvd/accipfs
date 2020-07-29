package main

import (
	"context"
	"fmt"
	"github.com/glvd/accipfs/client"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
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
	var byid bool
	connect := &cobra.Command{
		Use:   "connect",
		Short: "connect run",
		Long:  "connect a remote node",
		Run: func(cmd *cobra.Command, args []string) {
			config.Initialize()
			cfg := config.Global()
			client.InitGlobalClient(&cfg)
			ctx, cancelFunc := context.WithCancel(context.TODO())
			done := make(chan error)
			fmt.Printf("connect to %v\n", args)
			go func(c context.Context) {
				var err error
				defer func() {
					done <- err
				}()

				req := &core.NodeLinkReq{Addrs: args, ByID: byid, Names: args}
				resp, err := client.NodeLink(c, req)
				if err != nil {
					fmt.Printf("connect error: %v\n", err)
					return
				}
				fmt.Println("success:")
				for _, info := range resp.NodeInfos {
					fmt.Printf("connected to:%+v\n", info.ID)
				}
			}(ctx)
			sigs := make(chan os.Signal)
			signal.Notify(sigs, os.Interrupt)
			select {
			case <-sigs:
				cancelFunc()
			case v := <-done:
				if v != nil {
					panic(v)
				}
			}
			return
		},
	}
	connect.Flags().BoolVarP(&byid, "byid", "i", false, "connect to node by id")
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
			ctx, cancelFunc := context.WithCancel(context.TODO())
			done := make(chan error)
			go func(c context.Context) {
				var err error
				defer func() {
					done <- err
				}()

				list, err := client.NodeList(c, &core.NodeListReq{})
				if err != nil {
					panic(err)
				}
				for id := range list.Nodes {
					fmt.Printf("id:%v\n", id)
					//fmt.Printf("info:%v\n", info.JSON())
				}
			}(ctx)
			sigs := make(chan os.Signal)
			signal.Notify(sigs, os.Interrupt)
			select {
			case <-sigs:
				cancelFunc()
			case v := <-done:
				if v != nil {
					panic(v)
				}
			}
			return
		},
	}
	return peers
}
