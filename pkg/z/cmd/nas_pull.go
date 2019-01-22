package cmd

import (
	"context"
	"fmt"

	"github.com/cybozu-go/well"
	"github.com/mitsutaka/zcmd/nas"
	"github.com/spf13/cobra"
)

var nasPullCmdOpts struct {
	dryRun    bool
	syncPaths []string
}

// nasPullCmd represents the nas pull command
var nasPullCmd = &cobra.Command{
	Use:   "pull [PATH|all]",
	Short: "pull command pulls given PATH data to given PATH local directory",
	Long: `pull command pulls given PATH data to given PATH local directory.

-n option executes as dry-run.
all PATH pull all given paths in configuration file.`,
	Args: func(cmd *cobra.Command, args []string) error {
		cobra.MinimumNArgs(1)
		for _, path := range cfg.Nas.Pull.Sync {
			nasPullCmdOpts.syncPaths = append(nasPullCmdOpts.syncPaths, path.Name)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("%#v\n", nasPullCmdOpts.syncPaths)
		sync := nas.NewPull(&cfg.Nas.Pull, args, nasPullCmdOpts.dryRun)

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
