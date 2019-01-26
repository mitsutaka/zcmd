package cmd

import (
	"github.com/spf13/cobra"
)

// reposCmd represents the repos command
//nolint[gochecknoglobals]
var reposCmd = &cobra.Command{
	Use:   "repos",
	Short: "repos subcommand is operation for checked out git repositories.",
	Long:  `repos subcommand is operation for checked out local git repositories.`,
}

//nolint[gochecknoinits]
func init() {
	rootCmd.AddCommand(reposCmd)
}
