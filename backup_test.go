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
		name       string
		cfg        BackupConfig
		rsyncFlags string
		expected   []rsyncClient
	}{
		{
			name: "1-destination-3-sources",
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
			name: "1-destination-3-sources-rsync-flags",
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
			name: "1-rsync-destination-1-source",
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
		cfg := c.cfg
		rsyncFlags := c.rsyncFlags
		expected := c.expected

		t.Run(c.name, func(t *testing.T) {
			bk := NewBackup(&cfg, rsyncFlags)
			rcs, err := bk.generateCmd()

			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(rcs, expected) {
				t.Errorf("%#v != %#v", rcs, expected)
			}
		})
	}
}
