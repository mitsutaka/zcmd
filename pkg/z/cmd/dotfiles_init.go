package cmd

import (
	"context"
	"os"

	"github.com/cybozu-go/well"
	"github.com/mitsutaka/zcmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// dotfilesInitCmd represents the dotfiles init command
// nolint[gochecknoglobals]
var dotfilesInitCmd = &cobra.Command{
	Use:   "init",
	Short: "init setup dotfiles with your git dotfiles repository",
	Long:  `init setup dotfiles with your git dotfiles repository.`,
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		df, err := zcmd.NewDotFiler(&cfg.DotFiles)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		well.Go(func(ctx context.Context) error {
			return df.Init(ctx, args[0])
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
	dotfilesCmd.AddCommand(dotfilesInitCmd)
}
