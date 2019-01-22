package zcmd

import (
	"context"
	"fmt"
	"os"
	"runtime"
)

const (
	// CmdRsyncDarwin is default rsync absolute path for macOS
	cmdRsyncDarwin = "/usr/local/bin/rsync"
	// CmdRsyncLinux is default rsync absolute path for Linux
	cmdRsyncLinux = "/usr/bin/rsync"
	// OptDryRun is dry run option of rsync
	OptDryRun = "--dry-run"
)

var (
	// OptsRsync is default rsync options
	OptsRsync = []string{"-avxRP", "--stats", "--delete"}
	// sudoCmd is default sudo command
	sudoCmd = []string{"/usr/bin/sudo", "-E"}
)

// Rsync is rsync interface
type Rsync interface {
	// Do runs rsync command
	Do(ctx context.Context) error

	// GenerateCmd returns generated rsync commands
	GenerateCmd() (map[string][]string, error)
}

// GetRsyncCmd returns rsync command and arguments for each platform
func GetRsyncCmd() ([]string, error) {
	var cmdPrefix []string
	if os.Getuid() != 0 {
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
