package zcmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
)

const (
	backupPidFile = "/tmp/backup.pid"
	datePath      = "backup-0000-00-00-000000"
)

// Backup is client for backup
type Backup struct {
	destinations []string
	includes     []string
	excludes     []string
	rsyncFlags   string
}

// NewBackup returns Syncer
func NewBackup(cfg *BackupConfig, rsyncFlags string) Rsync {
	return &Backup{
		includes:     cfg.Includes,
		excludes:     cfg.Excludes,
		destinations: cfg.Destinations,
		rsyncFlags:   rsyncFlags,
	}
}

// Do is main backup process
func (b *Backup) Do(ctx context.Context) error {
	rcs, err := b.generateCmd()
	if err != nil {
		return err
	}

	return runRsyncCmd(ctx, "backup", backupPidFile, rcs)
}

// GenerateCmd generates rsync command
func (b *Backup) generateCmd() ([]rsyncClient, error) {
	var optsRsync = []string{"-avxRP", "--stats", "--delete"}

	cmdRsync, err := getRsyncCmd()
	if err != nil {
		return nil, err
	}

	cmdRsync = append(cmdRsync, optsRsync...)
	hostname, err := os.Hostname()

	if err != nil {
		return nil, err
	}

	i := 0
	cmds := make([]rsyncClient, len(b.includes)*len(b.destinations))

	for _, src := range b.includes {
		for _, dst := range b.destinations {
			var excludeFile *os.File
			if len(b.excludes) != 0 {
				excludeFile, err = ioutil.TempFile("", "backup")
				if err != nil {
					return nil, err
				}
				defer excludeFile.Close()

				for _, path := range b.excludes {
					_, err := excludeFile.WriteString(path + "\n")
					if err != nil {
						return nil, err
					}
				}
			}

			var cmd []string

			u, err := url.Parse(dst)

			if err != nil {
				return nil, err
			}

			u.Path = path.Join(u.Path, hostname, datePath)
			dst = u.String()

			cmd = append(cmd, cmdRsync...)

			if excludeFile != nil {
				cmd = append(cmd, fmt.Sprintf("--exclude-from=%s", excludeFile.Name()))
			}

			if len(b.rsyncFlags) != 0 {
				cmd = append(cmd, b.rsyncFlags)
			}

			cmd = append(cmd, src, dst)
			cmds[i] = rsyncClient{
				command:     cmd,
				excludeFile: excludeFile,
			}
			i++
		}
	}

	return cmds, nil
}
