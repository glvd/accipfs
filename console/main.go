//go:generate tar zxvf ../dhcrypto.tar.gz -C ../dhcrypto
package main

import (
	"fmt"
	"github.com/glvd/accipfs"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/general"
	"github.com/spf13/cobra"
	"path/filepath"
)

// APP ...
const APP = "accipfs"

// Version ...
const Version = "0.0.1"

var rootCmd = &cobra.Command{
	Use:   APP,
	Short: "accipfs is a very fast ipfs client",
	Long:  `accipfs`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func main() {
	path, err := filepath.Abs(accipfs.DefaultPath)
	if err != nil {
		path = general.CurrentDir()
	}
	config.WorkDir = path
	rootCmd.AddCommand(initCmd(), daemonCmd(), nodeCmd(), versionCmd(), tagCmd())
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
	rootCmd.Flags().StringVar(&accipfs.DefaultPath, "path", ".", "set work path")
}

func versionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of " + APP,
		Long:  `All software has versions. This is ` + APP + `'s`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(APP, "version:", Version)
		},
	}
	return cmd
}
