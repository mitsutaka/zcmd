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
	pidFile        = "/tmp/backup.pid"
	datePath       = "backup-0000-00-00-000000"
	cmdRsyncDarwin = "/usr/local/bin/rsync"
	cmdRsyncLinux  = "/usr/bin/rsync"
	optDryRun      = "--dry-run"
)

var (
	optsRsync = []string{"-avxRP", "--stats", "--delete"}
	sudoCmd   = []string{"/usr/bin/sudo", "-E"}
)

// Backup is client for backup
type Backup struct {
	destinations []string
	includes     []string
	excludes     []string
	dryRun       bool
}

// NewBackup returns Syncer
func NewBackup(cfg *BackupConfig, dryRun bool) *Backup {
	return &Backup{
		includes:     cfg.Includes,
		excludes:     cfg.Excludes,
		destinations: cfg.Destinations,
		dryRun:       dryRun,
	}
}

// Backup is main backup process
func (b *Backup) Do(ctx context.Context) {
	rsyncCmds, exclude, err := b.generateRsyncCmd()
	if err != nil {
		return
	}

	pid, err := os.Create(pidFile)
	if err != nil {
		return
	}
	defer func() {
		pid.Close()
		os.Remove(pid.Name())
		os.Remove(exclude)
	}()
	_, err = pid.WriteString(string(os.Getpid()))
	if err != nil {
		return
	}

	env := well.NewEnvironment(ctx)
	for _, rsyncCmd := range rsyncCmds {
		rsyncCmd := rsyncCmd
		env.Go(func(ctx context.Context) error {
			log.Printf("backup command: %#v\n", rsyncCmd)
			cmd := exec.Command(rsyncCmd[0], rsyncCmd[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		})
	}
	env.Stop()
	_ = env.Wait()
}

func (b *Backup) generateRsyncCmd() (map[string][]string, string, error) {
	var cmdRsync []string
	switch runtime.GOOS {
	case "linux":
		cmdRsync = []string{cmdRsyncLinux}
	case "darwin":
		cmdRsync = []string{cmdRsyncDarwin}
	default:
		return nil, "", fmt.Errorf("platform %s does not support", runtime.GOOS)
	}
	cmdRsync = append(cmdRsync, optsRsync...)

	var cmdPrefix []string
	if os.Getuid() != 0 {
		cmdPrefix = sudoCmd
	}

	f, err := ioutil.TempFile("", "backup")
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	optExclude := ""
	if b.excludes != nil {
		for _, path := range b.excludes {
			_, err := f.WriteString(path + "\n")
			if err != nil {
				return nil, "", err
			}
		}
		optExclude = fmt.Sprintf("--exclude-from=%s", f.Name())
	}

	cmds := make(map[string][]string)
	hostname, err := os.Hostname()
	if err != nil {
		return nil, "", err
	}
	for _, src := range b.includes {
		for _, dst := range b.destinations {
			var cmd []string
			dst := fmt.Sprintf("%s/%s/%s", dst, hostname, datePath)
			cmd = append(cmd, cmdPrefix...)
			cmd = append(cmd, cmdRsync...)
			if b.dryRun {
				cmd = append(cmd, optDryRun)
			}
			cmd = append(cmd, optExclude, src, dst)
			cmds[src] = cmd
		}
	}

	return cmds, f.Name(), nil
}
