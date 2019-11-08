package cmd

import (
	"github.com/spf13/cobra"
)

// completionCmd represents the nas command
//nolint[gochecknoglobals]
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate bash completion scripts",
	Long:  `Generate completion scripts`,
}

//nolint[gochecknoinits]
func init() {
	rootCmd.AddCommand(completionCmd)
}
