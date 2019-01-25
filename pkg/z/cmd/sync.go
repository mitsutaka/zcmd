package cmd

import (
	"github.com/spf13/cobra"
)

// syncCmd represents the nas command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync subcommand synchronizes data with the remote server",
	Long:  `sync subcommand synchronizes data with the remote server.`,
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
