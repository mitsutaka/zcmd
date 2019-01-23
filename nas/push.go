package nas

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/cybozu-go/well"
	"github.com/mitsutaka/zcmd"
)

// Push is client for nas push
type Push struct {
	argSyncs     []string
	cfgSyncs     *[]zcmd.SyncInfo
	excludeFiles []string
	dryRun       bool
}

// NewPush returns Syncer
func NewPush(sync *[]zcmd.SyncInfo, argSyncs []string, dryRun bool) zcmd.Rsync {
	return &Push{
		argSyncs: argSyncs,
		cfgSyncs: sync,
		dryRun:   dryRun,
	}
}

// Do is main pushing process
func (p *Push) Do(ctx context.Context) error {
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
			log.Printf("push started: %#v\n", rsyncCmd)
			cmd := exec.Command(rsyncCmd[0], rsyncCmd[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				return err
			}
			log.Printf("push finished: %#v\n", rsyncCmd)
			return nil
		})
	}
	env.Stop()
	return env.Wait()

}

// GenerateCmd generates rsync command
func (p *Push) GenerateCmd() (map[string][]string, error) {
	cmdRsync, err := zcmd.GetRsyncCmd()
	if err != nil {
		return nil, err
	}
	cmdRsync = append(cmdRsync, optsRsync...)

	targetSyncs := findTargetSyncs(*p.cfgSyncs, p.argSyncs)

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
		dst := sync.Destination
		src := sync.Source
		cmd = append(cmd, cmdRsync...)
		if p.dryRun {
			cmd = append(cmd, zcmd.OptDryRun)
		}
		if len(optExclude) != 0 {
			cmd = append(cmd, optExclude)
		}
		// Add "/" to sync all files in the source URL directory
		if !strings.HasSuffix(src, "/") {
			src = src + "/"
		}
		cmd = append(cmd, src, dst)
		cmds[sync.Name] = cmd
	}

	return cmds, nil
}
