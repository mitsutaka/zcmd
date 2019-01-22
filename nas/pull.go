package nas

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"runtime"

	"github.com/cybozu-go/well"
	"github.com/mitsutaka/zcmd"
)

var optsRsync = []string{"-avP", "--stats", "--delete", "--delete-excluded"}

// Pull is client for nas pull
type Pull struct {
	url          string
	sync         []zcmd.SyncInfo
	excludeFiles []string
	dryRun       bool
}

// NewPull returns Syncer
func NewPull(cfg *zcmd.NasPullConfig, dryRun bool) zcmd.Rsync {
	return &Pull{
		url:    cfg.URL,
		sync:   cfg.Sync,
		dryRun: dryRun,
	}
}

// Do is main pulling process
func (p *Pull) Do(ctx context.Context) error {
	rsyncCmds, err := p.GenerateCmd()
	if err != nil {
		return err
	}

	defer func() {
		for _, f := range p.excludeFiles {
			os.Remove(f)
		}
	}()

	env := well.NewEnvironment(ctx)
	for _, rsyncCmd := range rsyncCmds {
		rsyncCmd := rsyncCmd
		env.Go(func(ctx context.Context) error {
			log.Printf("sync started: %#v\n", rsyncCmd)
			//			cmd := exec.Command(rsyncCmd[0], rsyncCmd[1:]...)
			//			cmd.Stdout = os.Stdout
			//			cmd.Stderr = os.Stderr
			//			err := cmd.Run()
			//			if err != nil {
			//				return err
			//			}
			//			log.Printf("backup finished: %#v\n", rsyncCmd)
			return nil
		})
	}
	env.Stop()
	return env.Wait()

}

// GenerateCmd generates rsync command
func (p *Pull) GenerateCmd() (map[string][]string, error) {
	var cmdRsync []string
	switch runtime.GOOS {
	case "linux":
		cmdRsync = []string{zcmd.CmdRsyncLinux}
	case "darwin":
		cmdRsync = []string{zcmd.CmdRsyncDarwin}
	default:
		return nil, fmt.Errorf("platform %s does not support", runtime.GOOS)
	}
	cmdRsync = append(cmdRsync, optsRsync...)

	var cmdPrefix []string
	if os.Getuid() != 0 {
		cmdPrefix = zcmd.SudoCmd
	}

	cmds := make(map[string][]string)
	for _, sync := range p.sync {
		f, err := ioutil.TempFile("", sync.Name)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		optExclude := ""
		if sync.Excludes != nil {
			for _, exclude := range sync.Excludes {
				_, err := f.WriteString(fmt.Sprintf("*%s*\n", exclude))
				if err != nil {
					return nil, err
				}
			}
			optExclude = fmt.Sprintf("--exclude-from=%s", f.Name())
			p.excludeFiles = append(p.excludeFiles, f.Name())
		}

		var cmd []string
		u, err := url.Parse(p.url)
		if err != nil {
			return nil, err
		}
		u.Path = path.Join(u.Path, sync.Source)
		dst := sync.Destination
		cmd = append(cmd, cmdPrefix...)
		cmd = append(cmd, cmdRsync...)
		if p.dryRun {
			cmd = append(cmd, zcmd.OptDryRun)
		}
		if len(optExclude) != 0 {
			cmd = append(cmd, optExclude)
		}
		cmd = append(cmd, u.String()+"/", dst)
		cmds[sync.Name] = cmd
	}

	return cmds, nil
}
