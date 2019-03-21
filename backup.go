package zcmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/cybozu-go/well"
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
	dryRun       bool
}

// NewBackup returns Syncer
func NewBackup(cfg *BackupConfig, dryRun bool) Rsync {
	return &Backup{
		includes:     cfg.Includes,
		excludes:     cfg.Excludes,
		destinations: cfg.Destinations,
		dryRun:       dryRun,
	}
}

// Do is main backup process
func (b *Backup) Do(ctx context.Context) error {
	rcs, err := b.generateCmd()
	if err != nil {
		return err
	}

	pid, err := os.Create(backupPidFile)
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
			}).Info("backup started")

			cmd := well.CommandContext(ctx, rc.command[0], rc.command[1:]...)

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				log.WithFields(log.Fields{
					"command": strings.Join(rc.command, " "),
					"error":   err,
				}).Error("backup finished")
				return err
			}
			log.WithFields(log.Fields{
				"command": strings.Join(rc.command, " "),
			}).Info("backup finished")
			return nil
		})
	}
	env.Stop()
	return env.Wait()
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
			if b.dryRun {
				cmd = append(cmd, OptDryRun)
			}
			if excludeFile != nil {
				cmd = append(cmd, fmt.Sprintf("--exclude-from=%s", excludeFile.Name()))
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
