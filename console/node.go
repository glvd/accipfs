package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/glvd/accipfs/service"
	"github.com/gorilla/rpc/v2/json2"
	"github.com/spf13/cobra"
	"net/http"
)

func nodeCmd() *cobra.Command {
	nodeCmd := &cobra.Command{
		Use:   "node",
		Short: "node run",
		Long:  "node can operate to change the parameters of some nodes",
	}
	nodeCmd.AddCommand(nodeConnectCmd())
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
			message, err := json2.EncodeClientRequest("Accelerate.ID", &service.Empty{})
			if err != nil {
				panic(err)
			}
			url := fmt.Sprintf("http://localhost:%d/rpc", cfg.Port)
			resp, err := http.Post(url, "application/json", bytes.NewReader(message))
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
			reply := new(core.NodeInfo)
			err = json2.DecodeClientResponse(resp.Body, reply)
			if err != nil {
				panic(err)
			}
			message2, err := json2.EncodeClientRequest("Accelerate.Connect", reply)
			if err != nil {
				panic(err)
			}
			remoteURL := fmt.Sprintf("http://%s/rpc", addr)
			resp2, err := http.Post(remoteURL, "application/json", bytes.NewReader(message2))
			if err != nil {
				panic(err)
			}
			defer resp2.Body.Close()
			reply2 := new(bool)
			err = json2.DecodeClientResponse(resp.Body, reply2)
			if err != nil {
				panic(err)
			}
			if !(*reply2) {
				panic(errors.New("failed connect to remote"))
			}
			return
		},
	}
	connect.Flags().StringVar(&addr, "addr", "localhost", "set a remote address to connect")
	return connect
}
