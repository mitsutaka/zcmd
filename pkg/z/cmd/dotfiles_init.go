package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// dotfilesInitCmd represents the dotfiles init command
//nolint[gochecknoglobals]
var dotfilesInitCmd = &cobra.Command{
	Use:   "init",
	Short: "init setup dotfiles with your git dotfiles repository",
	Long:  `init setup dotfiles with your git dotfiles repository.`,
	Run: func(_ *cobra.Command, args []string) {
		os.Exit(0)
	},
}

//nolint[gochecknoinits]
func init() {
	dotfilesCmd.AddCommand(dotfilesInitCmd)
}
