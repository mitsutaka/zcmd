package zcmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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
	excludeFile  string
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
	rsyncCmds, err := b.GenerateCmd()
	if err != nil {
		return err
	}

	pid, err := os.Create(backupPidFile)
	if err != nil {
		return err
	}
	defer func() {
		os.Remove(pid.Name())
		if len(b.excludeFile) != 0 {
			os.Remove(b.excludeFile)
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
			log.WithFields(log.Fields{
				"command": strings.Join(rsyncCmd, " "),
			}).Info("backup started")

			cmd := exec.CommandContext(ctx, rsyncCmd[0], rsyncCmd[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				log.WithFields(log.Fields{
					"command": strings.Join(rsyncCmd, " "),
					"error":   err,
				}).Error("backup finished")
				return err
			}
			log.WithFields(log.Fields{
				"command": strings.Join(rsyncCmd, " "),
			}).Info("backup finished")
			return nil
		})
	}
	env.Stop()
	return env.Wait()
}

// GenerateCmd generates rsync command
func (b *Backup) GenerateCmd() (map[string][]string, error) {
	var optsRsync = []string{"-avxRP", "--stats", "--delete"}

	cmdRsync, err := GetRsyncCmd()
	if err != nil {
		return nil, err
	}
	cmdRsync = append(cmdRsync, optsRsync...)

	f, err := ioutil.TempFile("", "backup")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	optExclude := ""
	if b.excludes != nil {
		for _, path := range b.excludes {
			_, err := f.WriteString(path + "\n")
			if err != nil {
				return nil, err
			}
		}
		optExclude = fmt.Sprintf("--exclude-from=%s", f.Name())
		b.excludeFile = f.Name()
	}

	cmds := make(map[string][]string)
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	for _, src := range b.includes {
		for _, dst := range b.destinations {
			var cmd []string
			dst := fmt.Sprintf("%s/%s/%s", dst, hostname, datePath)
			cmd = append(cmd, cmdRsync...)
			if b.dryRun {
				cmd = append(cmd, OptDryRun)
			}
			cmd = append(cmd, optExclude, src, dst)
			cmds[src] = cmd
		}
	}

	return cmds, nil
}
