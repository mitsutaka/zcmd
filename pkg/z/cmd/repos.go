package cmd

import (
	"github.com/spf13/cobra"
)

// reposCmd represents the repos command
var reposCmd = &cobra.Command{
	Use:   "repos",
	Short: "repos subcommand is operation for checked out git repositories.",
	Long:  `repos subcommand is operation for checked out local git repositories.`,
}

func init() {
	rootCmd.AddCommand(reposCmd)
}
