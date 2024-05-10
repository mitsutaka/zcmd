package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// completionBashCmd represents the completion command
// nolint[gochecknoglobals]
var completionBashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Generate bash completion scripts",
	Long: `To load completion run

. <(z completion bash)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(z completion bash)`,
	Run: func(_ *cobra.Command, args []string) {
		err := rootCmd.GenBashCompletion(os.Stdout)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	},
}

// nolint[gochecknoinits]
func init() {
	completionCmd.AddCommand(completionBashCmd)
}
