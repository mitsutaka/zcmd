package cmd

import (
	"context"

	"github.com/cybozu-go/well"
	"github.com/mitsutaka/zcmd/sync"
	"github.com/spf13/cobra"
)

var syncPushCmdOpts struct {
	dryRun    bool
	syncPaths []string
}

// syncPushCmd represents the sync push command
var syncPushCmd = &cobra.Command{
	Use:   "push [-n] PATH",
	Short: "push command pushes given PATH local directory to given PATH in the remote server",
	Long: `push command pushes given PATH local directory to given PATH in the remote server.

-n option executes as dry-run.

When PATH is not given, all PATHs in configuration file will be synchronized.`,
	Args: func(cmd *cobra.Command, args []string) error {
		for _, arg := range args {
			for _, path := range cfg.Sync.Push {
				if arg == path.Name {
					syncPushCmdOpts.syncPaths = append(syncPushCmdOpts.syncPaths, path.Name)
				}
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		sync := sync.NewPush(&cfg.Sync.Push, args, syncPushCmdOpts.dryRun)

		well.Go(func(ctx context.Context) error {
			return sync.Do(ctx)
		})
		well.Stop()
		return well.Wait()
	},
}

func init() {
	syncPushCmd.Flags().BoolVarP(&syncPushCmdOpts.dryRun, "dry-run", "n", false, "dry run")
	syncCmd.AddCommand(syncPushCmd)
}
