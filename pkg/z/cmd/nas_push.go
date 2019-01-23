package cmd

import (
	"context"

	"github.com/cybozu-go/well"
	"github.com/mitsutaka/zcmd/nas"
	"github.com/spf13/cobra"
)

var nasPushCmdOpts struct {
	dryRun    bool
	syncPaths []string
}

// nasPushCmd represents the nas push command
var nasPushCmd = &cobra.Command{
	Use:   "push [-n] PATH",
	Short: "push command pushes given PATH local directory to given PATH in the remote server",
	Long: `push command pushes given PATH local directory to given PATH in the remote server.

-n option executes as dry-run.

When PATH is not given, all PATHs in configuration file will be synchronized.`,
	Args: func(cmd *cobra.Command, args []string) error {
		for _, arg := range args {
			for _, path := range cfg.Nas.Push {
				if arg == path.Name {
					nasPushCmdOpts.syncPaths = append(nasPushCmdOpts.syncPaths, path.Name)
				}
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		sync := nas.NewPush(&cfg.Nas.Push, args, nasPushCmdOpts.dryRun)

		well.Go(func(ctx context.Context) error {
			return sync.Do(ctx)
		})
		well.Stop()
		return well.Wait()
	},
}

func init() {
	nasPushCmd.Flags().BoolVarP(&nasPushCmdOpts.dryRun, "dry-run", "n", false, "dry run")
	nasCmd.AddCommand(nasPushCmd)
}
