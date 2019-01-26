package sync

import (
	"reflect"
	"testing"

	"github.com/mitsutaka/zcmd"
)

func testFindTargetSyncs(t *testing.T) {
	t.Parallel()

	cases := []struct {
		cfgs     []*zcmd.SyncInfo
		args     []string
		expected []*zcmd.SyncInfo
	}{
		{
			cfgs: []*zcmd.SyncInfo{
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
			expected: []*zcmd.SyncInfo{
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
			cfgs: []*zcmd.SyncInfo{
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
			expected: []*zcmd.SyncInfo{
				{
					Name:        "foo",
					Source:      "/foo",
					Destination: "/tmp/foo",
				},
			},
		},
		{
			cfgs: []*zcmd.SyncInfo{
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
			expected: []*zcmd.SyncInfo{
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
		cfgs     []*zcmd.SyncInfo
		args     []string
		dryRun   bool
		expected map[string][]string
	}{
		{
			cfgs: []*zcmd.SyncInfo{
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
			expected: map[string][]string{
				"foo": {"/usr/bin/sudo", "-E", "/usr/bin/rsync",
					"-avP", "--stats", "--delete", "--delete-excluded", "/foo/", "/tmp/foo"},
				"bar": {"/usr/bin/sudo", "-E", "/usr/bin/rsync",
					"-avP", "--stats", "--delete", "--delete-excluded", "/bar/", "/tmp/bar"},
			},
		},
		{
			cfgs: []*zcmd.SyncInfo{
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
			expected: map[string][]string{
				"foo": {"/usr/bin/sudo", "-E", "/usr/bin/rsync",
					"-avP", "--stats", "--delete", "--delete-excluded", "/foo/", "/tmp/foo"},
			},
		},
		{
			cfgs: []*zcmd.SyncInfo{
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
			expected: map[string][]string{
				"foo": {"/usr/bin/sudo", "-E", "/usr/bin/rsync",
					"-avP", "--stats", "--delete", "--delete-excluded", "/foo/", "/tmp/foo"},
				"bar": {"/usr/bin/sudo", "-E", "/usr/bin/rsync",
					"-avP", "--stats", "--delete", "--delete-excluded", "/bar/", "/tmp/bar"},
			},
		},
	}

	for _, c := range cases {
		sync := NewSync(c.cfgs, c.args, c.dryRun)
		cmds, err := sync.GenerateCmd()
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(cmds, c.expected) {
			t.Errorf("%#v != %#v", cmds, c.expected)
		}
	}
}

func TestSync(t *testing.T) {
	t.Run("findTargetSyncs", testFindTargetSyncs)
	t.Run("generateCmd", testGenerateCmd)
}
