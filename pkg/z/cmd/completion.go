package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
//nolint[gochecknoglobals]
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate bash completion scripts",
	Long: `To load completion run

. <(z completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(z completion)`,
	Run: func(_ *cobra.Command, args []string) {
		err := rootCmd.GenBashCompletion(os.Stdout)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	},
}

//nolint[gochecknoinits]
func init() {
	rootCmd.AddCommand(completionCmd)
}
