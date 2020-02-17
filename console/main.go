package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "accipfs",
	Short: "accipfs is a very fast ipfs client",
	Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// Execute ...
func Execute() {

}
func main() {
	fmt.Println("accipfs starting...")

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}

}
