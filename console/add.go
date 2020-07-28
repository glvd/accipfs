package main

import (
	"fmt"
	"github.com/glvd/accipfs/client"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/spf13/cobra"
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
			fmt.Println("add path", args[0])
			file, err := client.UploadFile(&core.UploadReq{
				Path: args[0],
			})
			if err != nil {
				panic(err)
			}
			add, err := client.Add(&core.AddReq{
				Hash: file.Hash,
			})
			if err != nil {
				panic(err)
			}
			fmt.Println("success", add.IsSuccess)
			fmt.Println("result hash:", add.Hash)

		},
	}
	cmd.Flags().StringVar(&path, "path", "", "set the file dirctory path to add")
	cmd.Flags().StringVar(&info, "info", "", "set the file info to load")
	return cmd
}
