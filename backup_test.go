package zcmd

import (
	"os"
	"reflect"
	"testing"
)

func TestBackup(t *testing.T) {
	t.Parallel()

	hostname, err := os.Hostname()
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		cfg      BackupConfig
		dryRun   bool
		expected []rsyncClient
	}{
		{
			cfg: BackupConfig{
				Destinations: []string{"/backup"},
				Includes:     []string{"/", "/boot", "/home"},
			},
			dryRun: true,
			expected: []rsyncClient{
				{
					command: []string{"/usr/bin/sudo", "-E", "/usr/bin/rsync",
						"-avxRP", "--stats", "--delete", "--dry-run", "/", "/backup/" + hostname + "/backup-0000-00-00-000000"},
					excludeFile: nil,
				},
				{
					command: []string{"/usr/bin/sudo", "-E", "/usr/bin/rsync",
						"-avxRP", "--stats", "--delete", "--dry-run", "/boot", "/backup/" + hostname + "/backup-0000-00-00-000000"},
					excludeFile: nil,
				},
				{
					command: []string{"/usr/bin/sudo", "-E", "/usr/bin/rsync",
						"-avxRP", "--stats", "--delete", "--dry-run", "/home", "/backup/" + hostname + "/backup-0000-00-00-000000"},
					excludeFile: nil,
				},
			},
		},
		{
			cfg: BackupConfig{
				Destinations: []string{"/backup"},
				Includes:     []string{"/", "/boot", "/home"},
			},
			dryRun: false,
			expected: []rsyncClient{
				{
					command: []string{"/usr/bin/sudo", "-E", "/usr/bin/rsync",
						"-avxRP", "--stats", "--delete", "/", "/backup/" + hostname + "/backup-0000-00-00-000000"},
					excludeFile: nil,
				},
				{
					command: []string{"/usr/bin/sudo", "-E", "/usr/bin/rsync",
						"-avxRP", "--stats", "--delete", "/boot", "/backup/" + hostname + "/backup-0000-00-00-000000"},
					excludeFile: nil,
				},
				{
					command: []string{"/usr/bin/sudo", "-E", "/usr/bin/rsync",
						"-avxRP", "--stats", "--delete", "/home", "/backup/" + hostname + "/backup-0000-00-00-000000"},
					excludeFile: nil,
				},
			},
		},
	}

	for _, c := range cases {
		bk := NewBackup(&c.cfg, c.dryRun)
		rcs, err := bk.generateCmd()
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(rcs, c.expected) {
			t.Errorf("%#v != %#v", rcs, c.expected)
		}
	}
}
