package main

import (
	"fmt"
	"github.com/glvd/accipfs/client"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/spf13/cobra"
)

func pinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pin",
		Short: "pin a video to local",
		Long:  "pin a video to local for sharing",
		Run: func(cmd *cobra.Command, args []string) {
			config.Initialize()
			cfg := config.Global()
			client.InitGlobalClient(&cfg)

			pins, err := client.Pins(&core.DataStoreReq{})
			if err != nil {
				return
			}
			fmt.Println("show pin list:")
			for _, v := range pins.Pins {
				fmt.Println(v)
			}
		},
	}
	return cmd
}

func pinHashCmd() *cobra.Command {
	cmd := &cobra.Command{}
	return cmd
}
