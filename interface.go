package zcmd

import "context"

const (
	cmdRsyncDarwin = "/usr/local/bin/rsync"
	cmdRsyncLinux  = "/usr/bin/rsync"
	optDryRun      = "--dry-run"
)

var (
	optsRsync = []string{"-avxRP", "--stats", "--delete"}
	sudoCmd   = []string{"/usr/bin/sudo", "-E"}
)

// Rsync is rsync interface
type Rsync interface {
	// Do runs rsync command
	Do(ctx context.Context) error

	// GenerateCmd returns generated rsync commands, exclude file name
	GenerateCmd() (map[string][]string, string, error)
}
