package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// reposUpdateCmd represents the repos update command
var reposUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "update fetches and checkouts from remote master branch",
	Long:  `update fetches and checkouts from remote master branch`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("repos update called")
	},
}

func init() {
	reposCmd.AddCommand(reposUpdateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reposCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reposCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
