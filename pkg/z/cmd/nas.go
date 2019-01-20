package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// nasCmd represents the nas command
var nasCmd = &cobra.Command{
	Use:   "nas",
	Short: "nas subcommand synchronizes data with the remote server",
	Long:  `nas subcommand synchronizes data with the remote server.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("nas called")
	},
}

func init() {
	rootCmd.AddCommand(nasCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nasCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nasCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
