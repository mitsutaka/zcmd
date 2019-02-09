package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// dotfilesPullCmd represents the dotfiles init command
//nolint[gochecknoglobals]
var dotfilesPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "pull download latest changes and create symlinks",
	Long:  `pull download latest changes and create symlinks.`,
	Run: func(_ *cobra.Command, args []string) {
		os.Exit(0)
	},
}

//nolint[gochecknoinits]
func init() {
	dotfilesCmd.AddCommand(dotfilesPullCmd)
}
