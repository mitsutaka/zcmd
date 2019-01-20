package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// nasPullCmd represents the nas pull command
var nasPullCmd = &cobra.Command{
	Use:   "pull [-n] PATH",
	Short: "pull command pulls given PATH data to given PATH local directory",
	Long: `pull command pulls given PATH data to given PATH local directory.

-n option executes as dry-run.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("nas pull called")
	},
}

func init() {
	nasCmd.AddCommand(nasPullCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nasCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nasCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
