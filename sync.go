package zcmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cybozu-go/well"
	log "github.com/sirupsen/logrus"
)

// Sync is client for sync pull
type Sync struct {
	argSyncs     []string
	cfgSyncs     []*SyncInfo
	excludeFiles []string
	dryRun       bool
}

// NewSync returns Syncer
func NewSync(sync []*SyncInfo, argSyncs []string, dryRun bool) Rsync {
	return &Sync{
		argSyncs: argSyncs,
		cfgSyncs: sync,
		dryRun:   dryRun,
	}
}

// Do is main pulling process
func (s *Sync) Do(ctx context.Context) error {
	syncPidFile := "/tmp/sync.pid"

	rsyncCmds, err := s.GenerateCmd()
	if err != nil {
		return err
	}

	pid, err := os.Create(syncPidFile)
	if err != nil {
		return err
	}
	defer func() {
		os.Remove(pid.Name())
		for _, f := range s.excludeFiles {
			os.Remove(f)
		}
	}()
	_, err = pid.WriteString(strconv.Itoa(os.Getpid()))
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
			log.WithFields(log.Fields{
				"command": strings.Join(rsyncCmd, " "),
			}).Info("sync started")

			cmd := exec.CommandContext(ctx, rsyncCmd[0], rsyncCmd[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				log.WithFields(log.Fields{
					"command": strings.Join(rsyncCmd, " "),
					"error":   err,
				}).Error("sync finished")
				return err
			}
			log.WithFields(log.Fields{
				"command": strings.Join(rsyncCmd, " "),
			}).Info("sync finished")
			return nil
		})
	}
	env.Stop()
	return env.Wait()
}

// GenerateCmd generates rsync command
func (s *Sync) GenerateCmd() (map[string][]string, error) {
	var optsRsync = []string{"-avP", "--stats", "--delete", "--delete-excluded"}

	cmdRsync, err := GetRsyncCmd()
	if err != nil {
		return nil, err
	}
	cmdRsync = append(cmdRsync, optsRsync...)

	targetSyncs := findTargetSyncs(s.cfgSyncs, s.argSyncs)

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
			s.excludeFiles = append(s.excludeFiles, f.Name())
		}

		var cmd []string
		src := sync.Source
		dst := sync.Destination
		cmd = append(cmd, cmdRsync...)
		if s.dryRun {
			cmd = append(cmd, OptDryRun)
		}
		if len(optExclude) != 0 {
			cmd = append(cmd, optExclude)
		}
		// Add "/" to sync all files in the source URL directory
		if !strings.HasSuffix(src, "/") {
			src += "/"
		}
		cmd = append(cmd, src, dst)
		cmds[sync.Name] = cmd
	}

	return cmds, nil
}

func findTargetSyncs(cfgs []*SyncInfo, args []string) []*SyncInfo {
	if len(args) == 0 {
		// Sync all paths
		return cfgs
	}

	targetCfgs := make([]*SyncInfo, 0)

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
