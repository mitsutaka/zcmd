package cmd

import (
	"fmt"
	"os"

	"github.com/cybozu-go/log"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/mitsutaka/zcmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var cfg *zcmd.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "z",
	Short:        "mitZ's command utilities",
	Long:         `Mitz's personal command line utilities.`,
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.ErrorExit(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.z.yaml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".z" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".z")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Printf("config file Unmarshal error: %v\n", err)
		os.Exit(1)
	}
}
