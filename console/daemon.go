package main

import (
	"github.com/spf13/cobra"
)

func daemonRun() *cobra.Command {
	return &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			//c := cron.New(cron.WithSeconds())
			//c.AddJob("", nil)
		},
	}
}
