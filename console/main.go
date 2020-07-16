//go:generate tar zxvf ../dhcrypto.tar.gz -C ../dhcrypto
package main

import (
	"fmt"
	"github.com/glvd/accipfs"
	"github.com/glvd/accipfs/basis"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/log"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
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
		path = basis.CurrentDir()
	}
	config.WorkDir = path

	rootCmd.AddCommand(initCmd(), daemonCmd(), idCmd(), nodeCmd(), versionCmd(), tagCmd(), pinCmd(), addCmd(), accountCmd())
	rootCmd.PersistentFlags().StringVar(&accipfs.DefaultPath, "path", ".", "set work path")

	rootCmd.PersistentFlags().StringVar(&log.Output, "log-output", "stdout", "set the output log name")
	rootCmd.PersistentFlags().StringVar(&log.Level, "log-level", "info", "set the log level(info,debug,warn,error,dpanic,panic,fatal)")
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		fmt.Println("log init")
		log.InitLog()
	}
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
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

func waitingInterruptSignal() {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, os.Interrupt)
	<-sigs
}
