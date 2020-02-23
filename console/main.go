package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "accipfs",
	Short: "accipfs is a very fast ipfs client",
	Long:  `accipfs`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func main() {
	fmt.Println("accipfs starting...")

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}

}
