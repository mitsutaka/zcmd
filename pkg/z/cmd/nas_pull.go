package cmd

import (
	"context"

	"github.com/cybozu-go/well"
	"github.com/mitsutaka/zcmd/nas"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		sync := nas.NewPull(&cfg.Nas.Pull, nasPullCmdOpts.dryRun)

		well.Go(func(ctx context.Context) error {
			return sync.Do(ctx)
		})
		well.Stop()
		return well.Wait()
	},
}

func init() {
	nasPullCmd.Flags().BoolVarP(&nasPullCmdOpts.dryRun, "dry-run", "n", false, "dry run")
	nasCmd.AddCommand(nasPullCmd)
}
