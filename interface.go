package zcmd

import "context"

const (
	// CmdRsyncDarwin is default rsync absolute path for macOS
	CmdRsyncDarwin = "/usr/local/bin/rsync"
	// CmdRsyncLinux is default rsync absolute path for Linux
	CmdRsyncLinux = "/usr/bin/rsync"
	// OptDryRun is dry run option of rsync
	OptDryRun = "--dry-run"
)

var (
	// OptsRsync is default rsync options
	OptsRsync = []string{"-avxRP", "--stats", "--delete"}
	// SudoCmd is default sudo command
	SudoCmd = []string{"/usr/bin/sudo", "-E"}
)

// Rsync is rsync interface
type Rsync interface {
	// Do runs rsync command
	Do(ctx context.Context) error

	// GenerateCmd returns generated rsync commands
	GenerateCmd() (map[string][]string, error)
}
