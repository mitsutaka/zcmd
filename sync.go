package zcmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const syncPidFile = "/tmp/sync.pid"

// Sync is client for sync pull
type Sync struct {
	argSyncs   []string
	cfgSyncs   []SyncInfo
	rsyncFlags string
}

// NewSync returns Syncer
func NewSync(sync []SyncInfo, argSyncs []string, rsyncFlags string) Rsync {
	return &Sync{
		argSyncs:   argSyncs,
		cfgSyncs:   sync,
		rsyncFlags: rsyncFlags,
	}
}

// Do is main pulling process
func (s *Sync) Do(ctx context.Context) error {
	rcs, err := s.generateCmd("")
	if err != nil {
		return err
	}

	return runRsyncCmd(ctx, "sync", syncPidFile, rcs)
}

// GenerateCmd generates rsync command
func (s *Sync) generateCmd(_ string) ([]rsyncClient, error) {
	optsRsync := []string{"-avP", "--stats", "--delete", "--delete-excluded"}

	targetSyncs := findTargetSyncs(s.cfgSyncs, s.argSyncs)

	cmds := make([]rsyncClient, len(targetSyncs))

	for i, sync := range targetSyncs {
		var excludeFile *os.File

		if sync.Excludes != nil {
			var err error

			excludeFile, err = ioutil.TempFile("", sync.Name)
			if err != nil {
				return nil, err
			}
			defer excludeFile.Close()

			for _, exclude := range sync.Excludes {
				_, err = excludeFile.WriteString(fmt.Sprintf("*%s*\n", exclude))
				if err != nil {
					return nil, err
				}
			}
		}

		var cmd []string

		cmdRsync, err := getRsyncCmd(sync.DisableSudo)
		if err != nil {
			return nil, err
		}

		cmd = append(cmd, cmdRsync...)
		cmd = append(cmd, optsRsync...)

		src := sync.Source
		dst := sync.Destination

		if excludeFile != nil {
			cmd = append(cmd, fmt.Sprintf("--exclude-from=%s", excludeFile.Name()))
		}
		// Add "/" to sync all files in the source URL directory
		if !strings.HasSuffix(src, "/") {
			src += "/"
		}

		if len(s.rsyncFlags) != 0 {
			cmd = append(cmd, s.rsyncFlags)
		}

		cmd = append(cmd, src, dst)

		cmds[i] = rsyncClient{
			command:     cmd,
			excludeFile: excludeFile,
		}
	}

	return cmds, nil
}

func findTargetSyncs(cfgs []SyncInfo, args []string) []SyncInfo {
	if len(args) == 0 {
		// Sync all paths
		return cfgs
	}

	targetCfgs := make([]SyncInfo, 0)

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
