package cmd

import (
	"github.com/spf13/cobra"
)

// dotfilesCmd represents the dotfiles command
//nolint[gochecknoglobals]
var dotfilesCmd = &cobra.Command{
	Use:   "dotfiles",
	Short: "dotfiles subcommand is dotfiles management",
	Long:  `dotfiles subcommand is dotfiles management.`,
}

//nolint[gochecknoinits]
func init() {
	rootCmd.AddCommand(dotfilesCmd)
}
