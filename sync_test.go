package zcmd

import (
	"reflect"
	"testing"
)

func testFindTargetSyncs(t *testing.T) {
	t.Parallel()

	cases := []struct {
		cfgs     []*SyncInfo
		args     []string
		expected []*SyncInfo
	}{
		{
			cfgs: []*SyncInfo{
				{
					Name:        "foo",
					Source:      "/foo",
					Destination: "/tmp/foo",
				},
				{
					Name:        "bar",
					Source:      "/bar",
					Destination: "/tmp/bar",
					Excludes:    []string{"aaa", "bbb"},
				},
			},
			args: []string{},
			expected: []*SyncInfo{
				{
					Name:        "foo",
					Source:      "/foo",
					Destination: "/tmp/foo",
				},
				{
					Name:        "bar",
					Source:      "/bar",
					Destination: "/tmp/bar",
					Excludes:    []string{"aaa", "bbb"},
				},
			},
		},
		{
			cfgs: []*SyncInfo{
				{
					Name:        "foo",
					Source:      "/foo",
					Destination: "/tmp/foo",
				},
				{
					Name:        "bar",
					Source:      "/bar",
					Destination: "/tmp/bar",
					Excludes:    []string{"aaa", "bbb"},
				},
			},
			args: []string{"foo"},
			expected: []*SyncInfo{
				{
					Name:        "foo",
					Source:      "/foo",
					Destination: "/tmp/foo",
				},
			},
		},
		{
			cfgs: []*SyncInfo{
				{
					Name:        "foo",
					Source:      "/foo",
					Destination: "/tmp/foo",
				},
				{
					Name:        "bar",
					Source:      "/bar",
					Destination: "/tmp/bar",
					Excludes:    []string{"aaa", "bbb"},
				},
			},
			args: []string{"foo", "bar"},
			expected: []*SyncInfo{
				{
					Name:        "foo",
					Source:      "/foo",
					Destination: "/tmp/foo",
				},
				{
					Name:        "bar",
					Source:      "/bar",
					Destination: "/tmp/bar",
					Excludes:    []string{"aaa", "bbb"},
				},
			},
		},
	}

	for _, c := range cases {
		syncInfo := findTargetSyncs(c.cfgs, c.args)
		if !reflect.DeepEqual(syncInfo, c.expected) {
			t.Errorf("%#v != %#v", syncInfo, c.expected)
		}
	}
}

func testGenerateCmd(t *testing.T) {
	t.Parallel()

	cases := []struct {
		cfgs     []*SyncInfo
		args     []string
		dryRun   bool
		expected []rsyncClient
	}{
		{
			cfgs: []*SyncInfo{
				{
					Name:        "foo",
					Source:      "/foo",
					Destination: "/tmp/foo",
				},
				{
					Name:        "bar",
					Source:      "/bar",
					Destination: "/tmp/bar",
				},
			},
			args: []string{},
			expected: []rsyncClient{
				{
					command: []string{"/usr/bin/sudo", "-E", "/usr/bin/rsync",
						"-avP", "--stats", "--delete", "--delete-excluded", "/foo/", "/tmp/foo"},
					excludeFile: nil,
				},
				{
					command: []string{"/usr/bin/sudo", "-E", "/usr/bin/rsync",
						"-avP", "--stats", "--delete", "--delete-excluded", "/bar/", "/tmp/bar"},
					excludeFile: nil,
				},
			},
		},
		{
			cfgs: []*SyncInfo{
				{
					Name:        "foo",
					Source:      "/foo",
					Destination: "/tmp/foo",
				},
				{
					Name:        "bar",
					Source:      "/bar",
					Destination: "/tmp/bar",
				},
			},
			args: []string{"foo"},
			expected: []rsyncClient{
				{
					command: []string{"/usr/bin/sudo", "-E", "/usr/bin/rsync",
						"-avP", "--stats", "--delete", "--delete-excluded", "/foo/", "/tmp/foo"},
					excludeFile: nil,
				},
			},
		},
		{
			cfgs: []*SyncInfo{
				{
					Name:        "foo",
					Source:      "/foo",
					Destination: "/tmp/foo",
				},
				{
					Name:        "bar",
					Source:      "/bar",
					Destination: "/tmp/bar",
				},
			},
			args: []string{"foo", "bar"},
			expected: []rsyncClient{
				{
					command: []string{"/usr/bin/sudo", "-E", "/usr/bin/rsync",
						"-avP", "--stats", "--delete", "--delete-excluded", "/foo/", "/tmp/foo"},
					excludeFile: nil,
				},
				{
					command: []string{"/usr/bin/sudo", "-E", "/usr/bin/rsync",
						"-avP", "--stats", "--delete", "--delete-excluded", "/bar/", "/tmp/bar"},
					excludeFile: nil,
				},
			},
		},
		{
			cfgs: []*SyncInfo{
				{
					Name:        "foo",
					Source:      "rsync://localhost/foo",
					Destination: "/tmp/foo",
				},
			},
			args: []string{"foo"},
			expected: []rsyncClient{
				{
					command: []string{"/usr/bin/sudo", "-E", "/usr/bin/rsync",
						"-avP", "--stats", "--delete", "--delete-excluded", "rsync://localhost/foo/", "/tmp/foo"},
					excludeFile: nil,
				},
			},
		},
	}

	for _, c := range cases {
		sync := NewSync(c.cfgs, c.args, c.dryRun)
		rcs, err := sync.generateCmd()
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(rcs, c.expected) {
			t.Errorf("%#v != %#v", rcs, c.expected)
		}
	}
}

func TestSync(t *testing.T) {
	t.Run("findTargetSyncs", testFindTargetSyncs)
	t.Run("generateCmd", testGenerateCmd)
}
