package zcmd

import (
	"context"
	"log"
	"os"

	"github.com/cybozu-go/well"
)

const (
	pidFile  = "/tmp/backup.pid"
	datePath = "backup-0000-00-00-000000"
)

var (
	rsyncOpts = []string{"-avxRP", "--stats", "--delete"}
	sudoCmd   = []string{"/usr/bin/sudo", "-E"}
	pid       = os.Getpid()
)

// Backup is client for backup
type Backup struct {
	*BackupConfig
}

// NewBackup returns Syncer
func NewBackup(cfg *BackupConfig) *Backup {
	return &Backup{cfg}
}

// Backup is main backup process
func (b *Backup) Do(ctx context.Context) {
	env := well.NewEnvironment(ctx)
	for _, d := range b.Destinations {
		for _, i := range b.Includes {
			dest := d
			include := i

			env.Go(func(ctx context.Context) error {
				log.Printf("backup %s, %s\n", dest, include)
				return nil
			})
		}
	}
	env.Stop()
	_ = env.Wait()
}
