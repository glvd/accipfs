package main

import "github.com/spf13/cobra"

func tagCmd() *cobra.Command {
	return &cobra.Command{}
}

func tagListCmd() *cobra.Command {
	return &cobra.Command{}
}

func tagAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "add a source to this node",
		Long:  "add a source to this node for service some users by your self",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}
