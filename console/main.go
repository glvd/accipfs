//go:generate tar zxvf ../dhcrypto.tar.gz -C ../dhcrypto
package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/glvd/accipfs"
	"github.com/glvd/accipfs/basis"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/log"
	"github.com/spf13/cobra"
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
var _pprofIP = "0.0.0.0:6060"

func main() {
	go func() {
		if err := http.ListenAndServe(_pprofIP, nil); err != nil {
			fmt.Printf("start pprof failed on %s\n", _pprofIP)
		}
	}()

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

func sysRoutineRun(f func(ctx context.Context, ch chan<- error)) {
	ctx, cancelFunc := context.WithCancel(context.TODO())
	done := make(chan error)
	sigs := make(chan os.Signal)
	signal.Notify(sigs, os.Interrupt)
	if f != nil {
		f(ctx, done)
	}
	select {
	case <-sigs:
		cancelFunc()
	case v := <-done:
		if v != nil {
			cancelFunc()
			panic(v)
		}
	}
}
