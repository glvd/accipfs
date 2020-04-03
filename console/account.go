package main

import (
	"encoding/json"
	"fmt"
	"github.com/glvd/accipfs/account"
	"github.com/glvd/accipfs/config"
	"github.com/spf13/cobra"
	"io/ioutil"
)

func accountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account",
		Short: "Account info",
		Long:  "Account show the information with your account",
	}
	cmd.AddCommand(accountInfoCmd(), accountSaveCmd())
	return cmd
}

func accountInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Info shows account information",
		Long:  "Info shows your account information with json format",
		Run: func(cmd *cobra.Command, args []string) {
			config.Initialize()
			cfg := config.Global()
			loadAccount, err := account.LoadAccount(&cfg)
			if err != nil {
				panic(err)
			}
			indent, err := json.MarshalIndent(loadAccount, "", " ")
			if err != nil {
				panic(err)
			}
			fmt.Println("Account info:")
			fmt.Println(string(indent))
		},
	}
}

func accountSaveCmd() *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:   "save",
		Short: "save account to file",
		Long:  "save your account to a file for backup your account info",
		Run: func(cmd *cobra.Command, args []string) {
			config.Initialize()
			cfg := config.Global()
			loadAccount, err := account.LoadAccount(&cfg)
			if err != nil {
				panic(err)
			}
			indent, err := json.MarshalIndent(loadAccount, "", " ")
			if err != nil {
				panic(err)
			}
			err = ioutil.WriteFile(path, indent, 0755)
			if err != nil {
				return
			}
		},
	}
	cmd.Flags().StringVar(&path, "path", "account.json", "save account for backup")
	return cmd
}
