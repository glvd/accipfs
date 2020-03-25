package main

import (
	"encoding/json"
	"fmt"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/glvd/accipfs/general"
	"github.com/glvd/accipfs/service"
	"github.com/spf13/cobra"
)

func idCmd() *cobra.Command {

	return &cobra.Command{
		Use:   "id",
		Short: "id print out your id",
		Long:  `id print the id information output to screen`,
		Run: func(cmd *cobra.Command, args []string) {
			config.Initialize()
			cfg := config.Global()
			url := fmt.Sprintf("http://localhost:%d/rpc", cfg.Port)
			reply := new(core.NodeInfo)
			if err := general.RPCPost(url, "Accelerate.ID", &service.Empty{}, reply); err != nil {
				fmt.Println("local id error:", err.Error())
				return
			}
			indent, err := json.MarshalIndent(reply, "", " ")
			if err != nil {
				return
			}
			//output your id info to screen
			fmt.Println(string(indent))
		},
	}
}
