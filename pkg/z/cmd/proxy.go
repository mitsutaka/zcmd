package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// proxyCmd represents the proxy command
//nolint[gochecknoglobals]
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "make ssh forwarding",
	Long:  `make ssh forwarding in parallel.`,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(0)
	},
}

//nolint[gochecknoinits]
func init() {
	rootCmd.AddCommand(proxyCmd)
}
