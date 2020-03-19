package main

import (
	"github.com/spf13/cobra"
)

func nodeCmd() *cobra.Command {
	nodeCmd := &cobra.Command{
		Use:   "node",
		Short: "node run",
		Long:  "node can operate to change the parameters of some nodes",
	}
	nodeCmd.AddCommand(nodeConnectCmd())
}

func nodeConnectCmd() *cobra.Command {
	var Addr string
	connect := &cobra.Command{
		Use:   "connect",
		Short: "connect run",
		Long:  "connect a remote node",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	connect.Flags().StringVar(&Addr, "addr", "localhost", "set a remote address to connect")

}
