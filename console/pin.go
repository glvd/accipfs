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

func pinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pin",
		Short: "show some pin info",
		Long:  "show the video information of pins with local server",
	}
	cmd.AddCommand(pinLsCmd())
	return cmd
}

func pinLsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "pin list",
		Long:  "show all pins list",
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
				pins, err := client.DataStorePinLs(c, &core.DataStoreReq{})
				if err != nil {
					panic(err)
				}
				fmt.Println("show pin list:")
				for _, v := range pins.Pins {
					fmt.Println(v)
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
		},
	}
}
