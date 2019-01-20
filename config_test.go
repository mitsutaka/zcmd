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
nas:
  pull:
    source: url
    destinations:
      - name: foo
        path: /foo
      - name: bar
        path: /bar
        excludes:
          - aaa
          - bbb
`,
			expected: Config{
				Nas: NasConfig{
					Pull: NasPullConfig{
						Source: "url",
						Destinations: []PathInfo{
							{
								Name: "foo",
								Path: "/foo",
							},
							{
								Name:     "bar",
								Path:     "/bar",
								Excludes: []string{"aaa", "bbb"},
							},
						},
					},
				},
			},
		},
		{
			source: `
nas:
  push:
    sources:
      - name: foo
        path: /foo
      - name: bar
        path: /bar
        excludes:
          - aaa
          - bbb
    destination: url
`,
			expected: Config{
				Nas: NasConfig{
					Push: NasPushConfig{
						Sources: []PathInfo{
							{
								Name: "foo",
								Path: "/foo",
							},
							{
								Name:     "bar",
								Path:     "/bar",
								Excludes: []string{"aaa", "bbb"},
							},
						},
						Destination: "url",
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
  root: ~/repos
`,
			expected: Config{
				Repos: ReposConfig{
					Root: "~/repos",
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
