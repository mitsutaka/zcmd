package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var nasPullCmdOpts struct {
	dryRun bool
}

// nasPullCmd represents the nas pull command
var nasPullCmd = &cobra.Command{
	Use:   "pull [PATH|all]",
	Short: "pull command pulls given PATH data to given PATH local directory",
	Long: `pull command pulls given PATH data to given PATH local directory.

-n option executes as dry-run.
all PATH pull all given paths in configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("nas pull called")
	},
}

func init() {
	nasPullCmd.Flags().BoolVarP(&nasPullCmdOpts.dryRun, "dry-run", "n", false, "dry run")
	nasCmd.AddCommand(nasPullCmd)
}
