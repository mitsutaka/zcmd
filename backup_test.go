package zcmd

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestBackup(t *testing.T) {
	t.Parallel()

	now := time.Now().Format(time.RFC3339)
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
						"-avxRP", "--stats", "--delete", "/", "/backup/" + hostname + "/" + now},
					excludeFile: nil,
				},
				{
					command: []string{"/usr/bin/sudo", "-E", "/usr/bin/rsync",
						"-avxRP", "--stats", "--delete", "/boot", "/backup/" + hostname + "/" + now},
					excludeFile: nil,
				},
				{
					command: []string{"/usr/bin/sudo", "-E", "/usr/bin/rsync",
						"-avxRP", "--stats", "--delete", "/home", "/backup/" + hostname + "/" + now},
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
						"-avxRP", "--stats", "--delete", "-n", "/", "/backup/" + hostname + "/" + now},
					excludeFile: nil,
				},
				{
					command: []string{"/usr/bin/sudo", "-E", "/usr/bin/rsync",
						"-avxRP", "--stats", "--delete", "-n", "/boot", "/backup/" + hostname + "/" + now},
				},
				{
					command: []string{"/usr/bin/sudo", "-E", "/usr/bin/rsync",
						"-avxRP", "--stats", "--delete", "-n", "/home", "/backup/" + hostname + "/" + now},
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
						"-avxRP", "--stats", "--delete", "/", "rsync://localhost/backup/" + hostname + "/" + now},
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
			rcs, err := bk.generateCmd(now)

			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(rcs, expected) {
				t.Errorf("%#v != %#v", rcs, expected)
			}
		})
	}
}
