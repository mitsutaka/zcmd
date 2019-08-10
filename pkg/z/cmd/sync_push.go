package cmd

import (
	"context"
	"os"

	"github.com/cybozu-go/well"
	"github.com/mitsutaka/zcmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//nolint[dupl]

//nolint[gochecknoglobals]
var syncPushCmdOpts struct {
	dryRun    bool
	syncPaths []string
}

// syncPushCmd represents the sync push command
//nolint[gochecknoglobals]
var syncPushCmd = &cobra.Command{
	Use:   "push [-n] PATH",
	Short: "push command pushes given PATH local directory to given PATH in the remote server",
	Long: `push command pushes given PATH local directory to given PATH in the remote server.

-n option executes as dry-run.

When PATH is not given, all PATHs in configuration file will be synchronized.`,
	Args: func(_ *cobra.Command, args []string) error {
		for _, arg := range args {
			for _, path := range cfg.Sync.Push {
				if arg == path.Name {
					syncPushCmdOpts.syncPaths = append(syncPushCmdOpts.syncPaths, path.Name)
				}
			}
		}
		return nil
	},
	Run: func(_ *cobra.Command, args []string) {
		sync := zcmd.NewSync(cfg.Sync.Push, args, syncPushCmdOpts.dryRun)

		well.Go(func(ctx context.Context) error {
			return sync.Do(ctx)
		})
		well.Stop()
		err := well.Wait()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		os.Exit(0)
	},
}

//nolint[gochecknoinits]
func init() {
	syncPushCmd.Flags().BoolVarP(&syncPushCmdOpts.dryRun, "dry-run", "n", false, "dry run")
	syncCmd.AddCommand(syncPushCmd)
}
