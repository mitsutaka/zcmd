package nas

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"

	"github.com/cybozu-go/well"
	"github.com/mitsutaka/zcmd"
)

// Pull is client for nas pull
type Pull struct {
	argSyncs     []string
	url          string
	cfgSyncs     []zcmd.SyncInfo
	excludeFiles []string
	dryRun       bool
}

// NewPull returns Syncer
func NewPull(cfg *zcmd.NasPullConfig, argSyncs []string, dryRun bool) zcmd.Rsync {
	return &Pull{
		argSyncs: argSyncs,
		url:      cfg.URL,
		cfgSyncs: cfg.Sync,
		dryRun:   dryRun,
	}
}

// Do is main pulling process
func (p *Pull) Do(ctx context.Context) error {
	rsyncCmds, err := p.GenerateCmd()
	if err != nil {
		return err
	}

	pid, err := os.Create(nasPidFile)
	if err != nil {
		return err
	}
	defer func() {
		os.Remove(pid.Name())
		for _, f := range p.excludeFiles {
			os.Remove(f)
		}
	}()
	_, err = pid.WriteString(string(os.Getpid()))
	if err != nil {
		return err
	}
	err = pid.Close()
	if err != nil {
		return err
	}

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

func findTargetSyncs(cfgs []zcmd.SyncInfo, args []string) []zcmd.SyncInfo {
	if len(args) == 0 {
		return nil
	}
	if args[0] == syncAllPath {
		// Sync all paths
		return cfgs
	}

	targetCfgs := make([]zcmd.SyncInfo, 0)

	for _, cfg := range cfgs {
		for _, arg := range args {
			if cfg.Name == arg {
				targetCfgs = append(targetCfgs, cfg)
				break
			}
		}
	}
	return targetCfgs
}

// GenerateCmd generates rsync command
func (p *Pull) GenerateCmd() (map[string][]string, error) {
	cmdRsync, err := zcmd.GetRsyncCmd()
	if err != nil {
		return nil, err
	}
	cmdRsync = append(cmdRsync, optsRsync...)

	targetSyncs := findTargetSyncs(p.cfgSyncs, p.argSyncs)

	cmds := make(map[string][]string)
	for _, sync := range targetSyncs {
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
		cmd = append(cmd, cmdRsync...)
		if p.dryRun {
			cmd = append(cmd, zcmd.OptDryRun)
		}
		if len(optExclude) != 0 {
			cmd = append(cmd, optExclude)
		}
		// Add "/" to sync all files in the source URL directory
		cmd = append(cmd, u.String()+"/", dst)
		cmds[sync.Name] = cmd
	}

	return cmds, nil
}
