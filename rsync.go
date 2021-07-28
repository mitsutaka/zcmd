package zcmd

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/cybozu-go/well"
	log "github.com/sirupsen/logrus"
)

const (
	// CmdRsyncDarwin is default rsync absolute path for macOS
	cmdRsyncDarwin = "/usr/local/bin/rsync"
	// CmdRsyncLinux is default rsync absolute path for Linux
	cmdRsyncLinux = "/usr/bin/rsync"
)

//nolint[gochecknoglobals]
var (
	// sudoCmd is default sudo command
	sudoCmd = []string{"/usr/bin/sudo", "-E"}
)

type rsyncClient struct {
	command     []string
	excludeFile *os.File
}

// Rsync is rsync interface
type Rsync interface {
	// Do runs rsync command
	Do(ctx context.Context) error

	// generateCmd returns generated rsync commands
	generateCmd(string) ([]rsyncClient, error)
}

// getRsyncCmd returns rsync command and arguments for each platform
func getRsyncCmd(disableSudo bool) ([]string, error) {
	var cmdPrefix []string
	if os.Getuid() != 0 && !disableSudo {
		cmdPrefix = sudoCmd
	}

	var cmdRsync []string

	switch runtime.GOOS {
	case "linux":
		cmdRsync = []string{cmdRsyncLinux}
	case "darwin":
		cmdRsync = []string{cmdRsyncDarwin}
	default:
		return nil, fmt.Errorf("platform %s does not support", runtime.GOOS)
	}

	var cmd []string
	cmd = append(cmd, cmdPrefix...)
	cmd = append(cmd, cmdRsync...)

	return cmd, nil
}

func runRsyncCmd(ctx context.Context, name, pidFile string, rcs []rsyncClient) error {
	pid, err := os.Create(pidFile)
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
			}).Info(name + " started")

			cmd := well.CommandContext(ctx, rc.command[0], rc.command[1:]...)

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				log.WithFields(log.Fields{
					"command": strings.Join(rc.command, " "),
					"error":   err,
				}).Error(name + " finished")
				return err
			}
			log.WithFields(log.Fields{
				"command": strings.Join(rc.command, " "),
			}).Info(name + " finished")
			return nil
		})
	}

	env.Stop()

	return env.Wait()
}
