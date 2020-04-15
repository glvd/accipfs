package main

import (
	"fmt"
	"github.com/glvd/accipfs/client"
	"github.com/glvd/accipfs/config"
	"github.com/spf13/cobra"
)

func pinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pin",
		Short: "pin a video to local",
		Long:  "pin a video to local for sharing",
		Run: func(cmd *cobra.Command, args []string) {
			config.Initialize()
			for _, no := range args {
				err := client.PinVideo(config.RPCAddr(), no)
				if err != nil {
					fmt.Printf("failed to pin (%s) with error(%v)\n", no, err.Error())
					return
				}
			}
		},
	}
	return cmd
}

func pinHashCmd() *cobra.Command {
	cmd := &cobra.Command{}
	return cmd
}
