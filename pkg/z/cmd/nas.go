package cmd

import (
	"github.com/spf13/cobra"
)

// nasCmd represents the nas command
var nasCmd = &cobra.Command{
	Use:   "nas",
	Short: "nas subcommand synchronizes data with the remote server",
	Long:  `nas subcommand synchronizes data with the remote server.`,
}

func init() {
	rootCmd.AddCommand(nasCmd)
}
