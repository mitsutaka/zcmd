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
var backupOpts struct {
	dryRun bool
}

// backupCmd represents the backup command
//nolint[gochecknoglobals]
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "backup local data to the remote server",
	Long:  `backup starts backup process local data to the remote server`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		bk := zcmd.NewBackup(&cfg.Backup, backupOpts.dryRun)

		well.Go(func(ctx context.Context) error {
			return bk.Do(ctx)
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
	backupCmd.Flags().BoolVarP(&backupOpts.dryRun, "dry-run", "n", false, "dry run")
	rootCmd.AddCommand(backupCmd)
}
