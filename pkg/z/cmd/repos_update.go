package cmd

import (
	"context"
	"os"
	"runtime"

	"github.com/cybozu-go/well"
	"github.com/mitsutaka/zcmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//nolint[gochecknoglobals]
var reposUpdateOpts struct {
	dryRun bool
	job    int
}

// reposUpdateCmd represents the repos update command
//nolint[gochecknoglobals]
var reposUpdateCmd = &cobra.Command{
	Use:   "repos-update",
	Short: "git clean, checkout and pull",
	Long:  `git clean, checkout and pull`,
	Args:  cobra.ExactArgs(0),
	Run: func(_ *cobra.Command, args []string) {
		upd, err := zcmd.NewUpdater(cfg.Repos.Root)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = upd.FindRepositories()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		if reposUpdateOpts.dryRun {
			os.Exit(0)
		}

		well.Go(func(ctx context.Context) error {
			return upd.Update(ctx, reposUpdateOpts.job)
		})
		well.Stop()
		err = well.Wait()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		os.Exit(0)
	},
}

//nolint[gochecknoinits]
func init() {
	reposUpdateCmd.Flags().BoolVarP(&reposUpdateOpts.dryRun, "dry-run", "n", false, "only show git repositories")
	reposUpdateCmd.Flags().IntVarP(&reposUpdateOpts.job, "job", "j", runtime.NumCPU(), "number of jobs")
	rootCmd.AddCommand(reposUpdateCmd)
}
