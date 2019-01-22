package nas

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	b "github.com/docker/docker/integration/build"
	"github.com/mitsutaka/zcmd"
)

// NasPull is client for nas pull
type NasPull struct {
	source       string
	destinations []zcmd.PathInfo
	dryRun       bool
}

// NewNasPull returns Syncer
func NewNasPull(cfg *zcmd.NasPullConfig, dryRun bool) zcmd.Rsync {
	return &NasPull{
		source:       cfg.Source,
		destinations: cfg.Destinations,
		dryRun:       dryRun,
	}
}

// Backup is main backup process
func (n *NasPull) Do(ctx context.Context) error {
	return nil
}

func (n *NasPull) GenerateCmd() (map[string][]string, string, error) {
	var cmdRsync []string
	switch runtime.GOOS {
	case "linux":
		cmdRsync = []string{cmdRsyncLinux}
	case "darwin":
		cmdRsync = []string{cmdRsyncDarwin}
	default:
		return nil, "", fmt.Errorf("platform %s does not support", runtime.GOOS)
	}
	cmdRsync = append(cmdRsync, optsRsync...)

	var cmdPrefix []string
	if os.Getuid() != 0 {
		cmdPrefix = sudoCmd
	}

	f, err := ioutil.TempFile("", "nasPull")
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	optExclude := ""
	if n.excludes != nil {
		for _, path := range n.destinations {
			_, err := f.WriteString(fmt.Sprintf("*%s*\n", path)
			if err != nil {
				return nil, "", err
			}
		}
		optExclude = fmt.Sprintf("--exclude-from=%s", f.Name())
	}

	cmds := make(map[string][]string)
	hostname, err := os.Hostname()
	if err != nil {
		return nil, "", err
	}
	for _, src := range b.includes {
		for _, dst := range b.destinations {
			var cmd []string
			dst := fmt.Sprintf("%s/%s/%s", dst, hostname, datePath)
			cmd = append(cmd, cmdPrefix...)
			cmd = append(cmd, cmdRsync...)
			if b.dryRun {
				cmd = append(cmd, optDryRun)
			}
			cmd = append(cmd, optExclude, src, dst)
			cmds[src] = cmd
		}
	}

	return cmds, f.Name(), nil
}
