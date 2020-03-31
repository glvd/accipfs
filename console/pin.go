package main

import "github.com/spf13/cobra"

func pinCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pin",
		Short: "pin a video to local",
		Long:  "pin a video to local for sharing",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}
