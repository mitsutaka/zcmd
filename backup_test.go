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
		cfg        BackupConfig
		rsyncFlags string
		expected   []rsyncClient
	}{
		{
			cfg: BackupConfig{
				Destinations: []string{"/backup"},
				Includes:     []string{"/", "/boot", "/home"},
			},
			rsyncFlags: "",
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
		{
			cfg: BackupConfig{
				Destinations: []string{"/backup"},
				Includes:     []string{"/", "/boot", "/home"},
			},
			rsyncFlags: "-n",
			expected: []rsyncClient{
				{
					command: []string{"/usr/bin/sudo", "-E", "/usr/bin/rsync",
						"-avxRP", "--stats", "--delete", "-n", "/", "/backup/" + hostname + "/backup-0000-00-00-000000"},
					excludeFile: nil,
				},
				{
					command: []string{"/usr/bin/sudo", "-E", "/usr/bin/rsync",
						"-avxRP", "--stats", "--delete", "-n", "/boot", "/backup/" + hostname + "/backup-0000-00-00-000000"},
					excludeFile: nil,
				},
				{
					command: []string{"/usr/bin/sudo", "-E", "/usr/bin/rsync",
						"-avxRP", "--stats", "--delete", "-n", "/home", "/backup/" + hostname + "/backup-0000-00-00-000000"},
					excludeFile: nil,
				},
			},
		},
		{
			cfg: BackupConfig{
				Destinations: []string{"rsync://localhost/backup"},
				Includes:     []string{"/"},
			},
			rsyncFlags: "",
			expected: []rsyncClient{
				{
					command: []string{"/usr/bin/sudo", "-E", "/usr/bin/rsync",
						"-avxRP", "--stats", "--delete", "/", "rsync://localhost/backup/" + hostname + "/backup-0000-00-00-000000"},
					excludeFile: nil,
				},
			},
		},
	}

	for _, c := range cases {
		bk := NewBackup(&c.cfg, c.rsyncFlags)
		rcs, err := bk.generateCmd()
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(rcs, c.expected) {
			t.Errorf("%#v != %#v", rcs, c.expected)
		}
	}
}
