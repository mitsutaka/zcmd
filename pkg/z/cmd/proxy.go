package cmd

import (
	"context"
	"os"

	"github.com/cybozu-go/well"
	"github.com/mitsutaka/zcmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// proxyCmd represents the proxy command
//nolint[gochecknoglobals]
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "make ssh forwarding",
	Long:  `make ssh forwarding in parallel.`,
	Run: func(cmd *cobra.Command, args []string) {
		proxy, err := zcmd.NewProxy(cfg.Proxy)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		well.Go(func(ctx context.Context) error {
			return proxy.Run(ctx)
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
	rootCmd.AddCommand(proxyCmd)
}
