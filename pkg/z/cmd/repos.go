package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// reposCmd represents the repos command
var reposCmd = &cobra.Command{
	Use:   "repos",
	Short: "repos subcommand is operation for checked out git repositories.",
	Long:  `repos subcommand is manuual operation for checked out local git repositories.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("repos called")
	},
}

func init() {
	rootCmd.AddCommand(reposCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reposCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reposCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
