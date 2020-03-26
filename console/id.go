package main

import (
	"encoding/json"
	"fmt"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/service"
	"github.com/goextension/log"
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
			id, err := service.ID(url)
			if err != nil {
				log.Errorw("local id", "error", err)
				return
			}
			indent, err := json.MarshalIndent(id, "", " ")
			if err != nil {
				log.Errorw("json mashal", "error", err)
				return
			}
			//output your id info to screen
			fmt.Println(string(indent))
		},
	}
}
