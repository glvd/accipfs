package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/glvd/accipfs/client"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

func idCmd() *cobra.Command {

	return &cobra.Command{
		Use:   "id",
		Short: "id print out your id",
		Long:  `id print the id information output to screen`,
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
				id, err := client.ID(c, &core.IDReq{})
				if err != nil {
					fmt.Printf("get local id failed error(%v)\n", err)
					return
				}
				indent, err := json.MarshalIndent(id, "", " ")
				if err != nil {
					fmt.Printf("json marshal failed error(%v)\n", err)
					return
				}
				//output your id info to screen
				fmt.Println(string(indent))
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
