package cmd

import (
	"context"

	"github.com/cybozu-go/well"
	"github.com/mitsutaka/zcmd"
	"github.com/spf13/cobra"
)

var backupOpts struct {
	dryRun bool
}

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "backup local data to the remote server",
	Long:  `backup starts backup process local data to the remote server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		bk := zcmd.NewBackup(&cfg.Backup)

		well.Go(func(ctx context.Context) error {
			bk.Do(ctx)
			return nil
		})
		well.Stop()
		return well.Wait()
	},
}

func init() {
	backupCmd.Flags().BoolVarP(&backupOpts.dryRun, "dry-run", "n", false, "only result of sync")
	rootCmd.AddCommand(backupCmd)
}
