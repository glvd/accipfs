package main

import "github.com/spf13/cobra"

func tagCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tag",
		Short: "tag contract",
		Long:  "tag contract manages all you videos",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}

func tagListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list videos to screen",
		Long:  "list and output the video number to screen",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
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
