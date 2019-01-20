package cmd

import (
	"context"

	"github.com/cybozu-go/well"
	"github.com/mitsutaka/zcmd"
	"github.com/spf13/cobra"
)

var reposUpdateOpts struct {
	dryRun bool
}

// reposUpdateCmd represents the repos update command
var reposUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "update fetches and checkouts from remote master branch",
	Long:  `update fetches and checkouts from remote master branch`,
	RunE: func(cmd *cobra.Command, args []string) error {
		upd, err := zcmd.NewUpdater(cfg.Repos.Root)
		if err != nil {
			return err
		}

		err = upd.FindRepositories()
		if err != nil {
			return err
		}

		if reposUpdateOpts.dryRun {
			return nil
		}

		well.Go(func(ctx context.Context) error {
			upd.FetchRepositories(ctx)
			upd.CheckoutRepositories(ctx)
			return nil
		})
		well.Stop()
		return well.Wait()
	},
}

func init() {
	reposUpdateCmd.Flags().BoolVarP(&reposUpdateOpts.dryRun, "dry-run", "n", false, "only show git repositories")
	reposCmd.AddCommand(reposUpdateCmd)
}
