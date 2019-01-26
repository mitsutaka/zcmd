package cmd

import (
	"context"

	"github.com/cybozu-go/well"
	"github.com/mitsutaka/zcmd/sync"
	"github.com/spf13/cobra"
)

var syncPullCmdOpts struct {
	dryRun    bool
	syncPaths []string
}

// syncPullCmd represents the sync pull command
var syncPullCmd = &cobra.Command{
	Use:   "pull [PATH]",
	Short: "pull command pulls given PATH data to given PATH local directory",
	Long: `pull command pulls given PATH data to given PATH local directory.

-n option executes as dry-run.

When PATH is not given, all PATHs in configuration file will be synchronized.`,
	Args: func(cmd *cobra.Command, args []string) error {
		for _, arg := range args {
			for _, path := range cfg.Sync.Pull {
				if arg == path.Name {
					syncPullCmdOpts.syncPaths = append(syncPullCmdOpts.syncPaths, path.Name)
				}
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		sync := sync.NewPull(&cfg.Sync.Pull, args, syncPullCmdOpts.dryRun)

		well.Go(func(ctx context.Context) error {
			return sync.Do(ctx)
		})
		well.Stop()
		return well.Wait()
	},
}

func init() {
	syncPullCmd.Flags().BoolVarP(&syncPullCmdOpts.dryRun, "dry-run", "n", false, "dry run")
	syncCmd.AddCommand(syncPullCmd)
}