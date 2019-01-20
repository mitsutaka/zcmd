package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// nasPushCmd represents the nas push command
var nasPushCmd = &cobra.Command{
	Use:   "push [-n] PATH",
	Short: "push command pushes given PATH local direcotry to given PATH in the remote server",
	Long: `push command pushes given PATH local direcotry to given PATH in the remote server.

-n option executes as dry-run.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("nas push called")
	},
}

func init() {
	nasCmd.AddCommand(nasPushCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nasCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nasCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
