package cmd

import (
	"context"
	"os"

	"github.com/cybozu-go/well"
	"github.com/mitsutaka/zcmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// dotfilesPullCmd represents the dotfiles init command
// nolint[gochecknoglobals]
var dotfilesPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "pull download latest changes and create symlinks",
	Long:  `pull download latest changes and create symlinks.`,
	Args:  cobra.ExactArgs(0),
	Run: func(_ *cobra.Command, args []string) {
		df, err := zcmd.NewDotFiler(&cfg.DotFiles)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		well.Go(func(ctx context.Context) error {
			return df.Pull(ctx)
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

// nolint[gochecknoinits]
func init() {
	dotfilesCmd.AddCommand(dotfilesPullCmd)
}
