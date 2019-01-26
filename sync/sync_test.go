package sync

import (
	"reflect"
	"testing"

	"github.com/mitsutaka/zcmd"
)

func TestBackup(t *testing.T) {
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
