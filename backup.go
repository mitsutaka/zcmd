package zcmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/cybozu-go/well"
)

const (
	pidFile  = "/tmp/backup.pid"
	datePath = "backup-0000-00-00-000000"
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

	pid, err := os.Create(pidFile)
	if err != nil {
		return err
	}
	defer func() {
		pid.Close()
		os.Remove(pid.Name())
		if len(b.excludeFile) != 0 {
			os.Remove(b.excludeFile)
		}
	}()
	_, err = pid.WriteString(string(os.Getpid()))
	if err != nil {
		return err
	}

	env := well.NewEnvironment(ctx)
	for _, rsyncCmd := range rsyncCmds {
		rsyncCmd := rsyncCmd
		env.Go(func(ctx context.Context) error {
			log.Printf("backup started: %#v\n", rsyncCmd)
			cmd := exec.Command(rsyncCmd[0], rsyncCmd[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				return err
			}
			log.Printf("backup finished: %#v\n", rsyncCmd)
			return nil
		})
	}
	env.Stop()
	return env.Wait()
}

// GenerateCmd generates rsync command
func (b *Backup) GenerateCmd() (map[string][]string, error) {
	var cmdRsync []string
	switch runtime.GOOS {
	case "linux":
		cmdRsync = []string{CmdRsyncLinux}
	case "darwin":
		cmdRsync = []string{CmdRsyncDarwin}
	default:
		return nil, fmt.Errorf("platform %s does not support", runtime.GOOS)
	}
	cmdRsync = append(cmdRsync, OptsRsync...)

	var cmdPrefix []string
	if os.Getuid() != 0 {
		cmdPrefix = SudoCmd
	}

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
			cmd = append(cmd, cmdPrefix...)
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
