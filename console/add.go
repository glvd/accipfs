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

func addCmd() *cobra.Command {
	var path string
	var info string
	cmd := &cobra.Command{
		Use:   "add",
		Short: "add a source to this node",
		Long:  "add a source to this node for service some users by your self",
		Run: func(cmd *cobra.Command, args []string) {
			//add a file with rule to accipf
			config.Initialize()
			cfg := config.Global()
			client.InitGlobalClient(&cfg)
			if len(args) <= 0 {
				return
			}
			ctx, cancelFunc := context.WithCancel(context.TODO())
			done := make(chan error)
			fmt.Println("add path", args[0])
			go func(c context.Context) {
				var err error
				defer func() {
					done <- err
				}()
				file, err := client.UploadFile(c, &core.UploadReq{
					Path: args[0],
				})
				if err != nil {
					return
				}
				add, err := client.Add(c, &core.NodeAddReq{
					Hash: file.Hash,
				})
				if err != nil {
					return
				}
				fmt.Println("success", add.IsSuccess)
				fmt.Println("result hash:", add.Hash)
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
	cmd.Flags().StringVar(&path, "path", "", "set the file dirctory path to add")
	cmd.Flags().StringVar(&info, "info", "", "set the file info to load")
	return cmd
}
