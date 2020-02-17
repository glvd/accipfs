package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

func initCmd() *cobra.Command {

}

func initConfig() (cfgFile string) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			panic(err)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	return
}
