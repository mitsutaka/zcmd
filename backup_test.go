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
		expected map[string][]string
	}{
		{
			cfg: BackupConfig{
				Destinations: []string{"/backup"},
				Includes:     []string{"/", "/boot", "/home"},
			},
			dryRun: true,
			expected: map[string][]string{
				"/": {"/usr/bin/sudo", "-E", "/usr/bin/rsync",
					"-avxRP", "--stats", "--delete", "--dry-run", "", "/", "/backup/" + hostname + "/backup-0000-00-00-000000"},
				"/boot": {"/usr/bin/sudo", "-E", "/usr/bin/rsync",
					"-avxRP", "--stats", "--delete", "--dry-run", "", "/boot", "/backup/" + hostname + "/backup-0000-00-00-000000"},
				"/home": {"/usr/bin/sudo", "-E", "/usr/bin/rsync",
					"-avxRP", "--stats", "--delete", "--dry-run", "", "/home", "/backup/" + hostname + "/backup-0000-00-00-000000"},
			},
		},
		{
			cfg: BackupConfig{
				Destinations: []string{"/backup"},
				Includes:     []string{"/", "/boot", "/home"},
			},
			dryRun: false,
			expected: map[string][]string{
				"/": {"/usr/bin/sudo", "-E", "/usr/bin/rsync",
					"-avxRP", "--stats", "--delete", "", "/", "/backup/" + hostname + "/backup-0000-00-00-000000"},
				"/boot": {"/usr/bin/sudo", "-E", "/usr/bin/rsync",
					"-avxRP", "--stats", "--delete", "", "/boot", "/backup/" + hostname + "/backup-0000-00-00-000000"},
				"/home": {"/usr/bin/sudo", "-E", "/usr/bin/rsync",
					"-avxRP", "--stats", "--delete", "", "/home", "/backup/" + hostname + "/backup-0000-00-00-000000"},
			},
		},
	}

	for _, c := range cases {
		bk := NewBackup(&c.cfg, c.dryRun)
		cmds, err := bk.GenerateCmd()
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(cmds, c.expected) {
			t.Errorf("%#v != %#v", cmds, c.expected)
		}
	}
}
