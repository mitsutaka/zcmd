package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// completionZshCmd represents the completion command
// nolint[gochecknoglobals]
var completionZshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Generate zsh completion scripts",
	Long: `To load completion run

. <(z completion zsh)

To configure your zsh shell to load completions for each session add to your zshrc

# ~/.zshrc or ~/.profile
. <(z completion zsh)`,
	Run: func(_ *cobra.Command, args []string) {
		err := rootCmd.GenZshCompletion(os.Stdout)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	},
}

// nolint[gochecknoinits]
func init() {
	completionCmd.AddCommand(completionZshCmd)
}
