package cmd

import (
	"context"
	"os"

	"github.com/cybozu-go/well"
	"github.com/mitsutaka/zcmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//nolint[gochecknoglobals]
var syncPullCmdOpts struct {
	dryRun    bool
	syncPaths []string
}

// syncPullCmd represents the sync pull command
//nolint[gochecknoglobals]
var syncPullCmd = &cobra.Command{
	Use:   "pull [PATH]",
	Short: "pull command pulls given PATH data to given PATH local directory",
	Long: `pull command pulls given PATH data to given PATH local directory.

-n option executes as dry-run.

When PATH is not given, all PATHs in configuration file will be synchronized.`,
	Args: func(_ *cobra.Command, args []string) error {
		for _, arg := range args {
			for _, path := range cfg.Sync.Pull {
				if arg == path.Name {
					syncPullCmdOpts.syncPaths = append(syncPullCmdOpts.syncPaths, path.Name)
				}
			}
		}
		return nil
	},
	Run: func(_ *cobra.Command, args []string) {
		sync := zcmd.NewSync(cfg.Sync.Pull, args, syncPullCmdOpts.dryRun)

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
	syncPullCmd.Flags().BoolVarP(&syncPullCmdOpts.dryRun, "dry-run", "n", false, "dry run")
	syncCmd.AddCommand(syncPullCmd)
}
