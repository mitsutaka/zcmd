package zcmd

import (
	"reflect"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	cases := []struct {
		source   string
		expected Config
	}{
		{
			source: `
sync:
  pull:
    - name: foo
      source: /foo
      destination: /tmp/foo
    - name: bar
      source: /bar
      destination: /tmp/bar
      excludes:
        - aaa
        - bbb
`,
			expected: Config{
				Sync: SyncConfig{
					Pull: []*SyncInfo{
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
			},
		},
		{
			source: `
sync:
  push:
    - name: foo
      source: /foo
      destination: /tmp/foo
    - name: bar
      source: /bar
      destination: /tmp/bar
      excludes:
        - aaa
        - bbb
`,
			expected: Config{
				Sync: SyncConfig{
					Push: []*SyncInfo{
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
			},
		},
		{
			source: `
backup:
  destinations:
    - /backup
  includes:
    - /
    - /boot
    - /home
  excludes:
    - foo
    - bar
`,
			expected: Config{
				Backup: BackupConfig{
					Destinations: []string{"/backup"},
					Includes:     []string{"/", "/boot", "/home"},
					Excludes:     []string{"foo", "bar"},
				},
			},
		},
		{
			source: `
repos:
  root: /repos
`,
			expected: Config{
				Repos: ReposConfig{
					Root: "/repos",
				},
			},
		},
	}

	for _, c := range cases {
		cfg := NewConfig()
		err := yaml.Unmarshal([]byte(c.source), &cfg)
		if err != nil {
			t.Error(err)
		} else if !reflect.DeepEqual(*cfg, c.expected) {
			t.Errorf("%#v != %#v", *cfg, c.expected)
		}
	}
}
