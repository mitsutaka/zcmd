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

const syncPidFile = "/tmp/sync.pid"

// Sync is client for sync pull
type Sync struct {
	argSyncs []string
	cfgSyncs []*SyncInfo
	dryRun   bool
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
	rcs, err := s.generateCmd()
	if err != nil {
		return err
	}

	pid, err := os.Create(syncPidFile)
	if err != nil {
		return err
	}
	defer os.Remove(pid.Name())
	_, err = pid.WriteString(strconv.Itoa(os.Getpid()))
	if err != nil {
		return err
	}
	err = pid.Close()
	if err != nil {
		return err
	}

	env := well.NewEnvironment(ctx)
	for _, rc := range rcs {
		rc := rc

		env.Go(func(ctx context.Context) error {
			defer func() {
				if rc.excludeFile != nil {
					os.Remove(rc.excludeFile.Name())
				}
			}()

			log.WithFields(log.Fields{
				"command": strings.Join(rc.command, " "),
			}).Info("sync started")

			cmd := exec.CommandContext(ctx, rc.command[0], rc.command[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				log.WithFields(log.Fields{
					"command": strings.Join(rc.command, " "),
					"error":   err,
				}).Error("sync finished")
				return err
			}
			log.WithFields(log.Fields{
				"command": strings.Join(rc.command, " "),
			}).Info("sync finished")
			return nil
		})
	}
	env.Stop()
	return env.Wait()
}

// GenerateCmd generates rsync command
func (s *Sync) generateCmd() ([]rsyncClient, error) {
	var optsRsync = []string{"-avP", "--stats", "--delete", "--delete-excluded"}

	cmdRsync, err := GetRsyncCmd()
	if err != nil {
		return nil, err
	}
	cmdRsync = append(cmdRsync, optsRsync...)

	targetSyncs := findTargetSyncs(s.cfgSyncs, s.argSyncs)

	cmds := make([]rsyncClient, len(targetSyncs))
	for i, sync := range targetSyncs {
		var excludeFile *os.File
		if sync.Excludes != nil {
			excludeFile, err = ioutil.TempFile("", sync.Name)
			if err != nil {
				return nil, err
			}
			defer excludeFile.Close()

			for _, exclude := range sync.Excludes {
				_, err := excludeFile.WriteString(fmt.Sprintf("*%s*\n", exclude))
				if err != nil {
					return nil, err
				}
			}
		}

		var cmd []string
		src := sync.Source
		dst := sync.Destination
		cmd = append(cmd, cmdRsync...)
		if s.dryRun {
			cmd = append(cmd, OptDryRun)
		}
		if excludeFile != nil {
			cmd = append(cmd, fmt.Sprintf("--exclude-from=%s", excludeFile.Name()))
		}
		// Add "/" to sync all files in the source URL directory
		if !strings.HasSuffix(src, "/") {
			src += "/"
		}
		cmd = append(cmd, src, dst)

		cmds[i] = rsyncClient{
			command:     cmd,
			excludeFile: excludeFile,
		}
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
