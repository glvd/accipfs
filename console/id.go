package main

import (
	"encoding/json"
	"fmt"
	"github.com/glvd/accipfs/client"
	"github.com/glvd/accipfs/config"
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
			id, err := client.ID(url)
			if err != nil {
				fmt.Printf("get local id failed error(%v)", err)
				return
			}
			indent, err := json.MarshalIndent(id, "", " ")
			if err != nil {
				fmt.Printf("json marshal failed error(%v)", err)
				return
			}
			//output your id info to screen
			fmt.Println(string(indent))
		},
	}
}
