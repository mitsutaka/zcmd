package zcmd

import (
	"reflect"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		source   string
		expected *Config
	}{
		{
			name: "sync-1",
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
			expected: &Config{
				Sync: SyncConfig{
					Pull: []SyncInfo{
						{
							Name:        "foo",
							Source:      "/foo",
							Destination: "/tmp/foo",
							DisableSudo: false,
						},
						{
							Name:        "bar",
							Source:      "/bar",
							Destination: "/tmp/bar",
							Excludes:    []string{"aaa", "bbb"},
							DisableSudo: false,
						},
					},
				},
				DotFiles: DotFilesConfig{
					Dir: defaultDotFilesDir,
				},
			},
		},
		{
			name: "sync-2",
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
      disable_sudo: true
`,
			expected: &Config{
				Sync: SyncConfig{
					Push: []SyncInfo{
						{
							Name:        "foo",
							Source:      "/foo",
							Destination: "/tmp/foo",
							DisableSudo: false,
						},
						{
							Name:        "bar",
							Source:      "/bar",
							Destination: "/tmp/bar",
							Excludes:    []string{"aaa", "bbb"},
							DisableSudo: true,
						},
					},
				},
				DotFiles: DotFilesConfig{
					Dir: defaultDotFilesDir,
				},
			},
		},
		{
			name: "backup",
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
			expected: &Config{
				Backup: BackupConfig{
					Destinations: []string{"/backup"},
					Includes:     []string{"/", "/boot", "/home"},
					Excludes:     []string{"foo", "bar"}},
				DotFiles: DotFilesConfig{
					Dir: defaultDotFilesDir,
				},
			},
		},
		{
			name: "repos",
			source: `
repos:
  root: /repos
`,
			expected: &Config{
				Repos: ReposConfig{
					Root: "/repos",
				},
				DotFiles: DotFilesConfig{
					Dir: defaultDotFilesDir,
				},
			},
		},
		{
			name: "dotfiles-1",
			source: `
dotfiles:
  dir: /home/mitz/.dotfiles
  hosts:
    - YOUR_HOSTNAME
  files:
    - bashrc
    - config/sway/config
    - spacemacs
    - ssh
`,
			expected: &Config{
				DotFiles: DotFilesConfig{
					Dir:   "/home/mitz/.dotfiles",
					Hosts: []string{"YOUR_HOSTNAME"},
					Files: []string{"bashrc", "config/sway/config", "spacemacs", "ssh"},
				},
			},
		},
		{
			name: "dotfiles-2",
			source: `
dotfiles:
  hosts:
    - YOUR_HOSTNAME
  files:
    - bashrc
    - config/sway/config
    - spacemacs
    - ssh
`,
			expected: &Config{
				DotFiles: DotFilesConfig{
					Dir:   defaultDotFilesDir,
					Hosts: []string{"YOUR_HOSTNAME"},
					Files: []string{"bashrc", "config/sway/config", "spacemacs", "ssh"},
				},
			},
		},
		{
			name: "proxy",
			source: `
proxy:
  - name: testforward1
    user: ubuntu
    address: remotehost1
    privateKey: /home/mitz/.ssh/id_rsa
    forward:
      # Local forwarding
      - type: local
        # default bindAddress is *
        bindAddress: localhost
        bindPort: 13128
        remoteAddress: localhost
        remotePort: 3128
      # Dynamic forwarding for SOCK4, 5
      - type: dynamic
        bindAddress: localhost
        bindPort: 1080
  - name: testforward2
    user: admin
    address: remotehost2
    privateKey: /home/mitz/.ssh/id_ecdsa
    port: 10000
    forward:
      # Remote forwarding
      - type: remote
        bindAddress: localhost
        bindPort: 9000
        remoteAddress: localhost
        remotePort: 3000
`,
			expected: &Config{
				Proxy: []ProxyConfig{
					{
						Name:       "testforward1",
						User:       "ubuntu",
						Address:    "remotehost1",
						PrivateKey: "/home/mitz/.ssh/id_rsa",
						Port:       DefaultProxyPort,
						Forward: []ProxyForwardConfig{
							{
								Type:          LocalForward,
								BindAddress:   "localhost",
								BindPort:      13128,
								RemoteAddress: "localhost",
								RemotePort:    3128,
							},
							{
								Type:          DynamicForward,
								BindAddress:   "localhost",
								BindPort:      1080,
								RemoteAddress: "",
								RemotePort:    0,
							},
						},
					},
					{
						Name:       "testforward2",
						User:       "admin",
						Address:    "remotehost2",
						PrivateKey: "/home/mitz/.ssh/id_ecdsa",
						Port:       10000,
						Forward: []ProxyForwardConfig{
							{
								Type:          RemoteForward,
								BindAddress:   "localhost",
								BindPort:      9000,
								RemoteAddress: "localhost",
								RemotePort:    3000,
							},
						},
					},
				},
				DotFiles: DotFilesConfig{
					Dir: defaultDotFilesDir,
				},
			},
		},
	}

	for _, c := range cases {
		source := c.source
		expected := c.expected

		t.Run(c.name, func(t *testing.T) {
			cfg, err := NewConfig(source)
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(cfg, expected) {
				t.Errorf(pretty.Compare(cfg, expected))
			}
		})
	}
}
