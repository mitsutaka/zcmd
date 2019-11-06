package zcmd

import (
	"os"
	"path/filepath"
	"strings"

	"sigs.k8s.io/yaml"

	"github.com/mitchellh/go-homedir"
	"github.com/tcnksm/go-input"
)

type Config struct {
	Sync     SyncConfig     `json:"sync,omitempty"`
	Backup   BackupConfig   `json:"backup,omitempty"`
	Repos    ReposConfig    `json:"repos,omitempty"`
	DotFiles DotFilesConfig `json:"dotfiles,omitempty"`
	Proxy    []ProxyConfig  `json:"proxy,omitempty"`
}

// SyncConfig is sync: config
type SyncConfig struct {
	Pull []SyncInfo `json:"pull,omitempty"`
	Push []SyncInfo `json:"push,omitempty"`
}

// SyncInfo is path information for synchronize directories
type SyncInfo struct {
	Name        string   `json:"name"`
	Source      string   `json:"source"`
	Destination string   `json:"destination"`
	Excludes    []string `json:"excludes,omitempty"`
	DisableSudo bool     `json:"disable_sudo,omitempty"`
}

// BackupConfig is backup: config
type BackupConfig struct {
	Destinations []string `json:"destinations"`
	Includes     []string `json:"includes"`
	Excludes     []string `json:"excludes,omitempty"`
}

// ReposConfig is repos: config
type ReposConfig struct {
	Root string `json:"root"`
}

var (
	// nolint[gochecknoglobals]
	// defaultDotFilesDir is default dotfiles directory path
	defaultDotFilesDir = filepath.Join(os.Getenv("HOME"), ".zdotfiles")
)

// DotFilesConfig is dotfiles: config
type DotFilesConfig struct {
	Dir   string   `json:"dir,omitempty"`
	Hosts []string `json:"hosts"`
	Files []string `json:"files"`
}

// ProxyConfig is proxy: config
type ProxyConfig struct {
	Name       string               `json:"name"`
	User       string               `json:"user,omitempty"`
	Address    string               `json:"address"`
	Port       int                  `json:"port,omitempty"`
	PrivateKey string               `json:"privateKey"`
	Forward    []ProxyForwardConfig `json:"forward"`
}

// ProxyForwardType is ssh forwarding type
type ProxyForwardType string

const (
	// DefaultProxyPort is default ssh port
	DefaultProxyPort int = 22
	// LocalForward is local forwarding
	LocalForward ProxyForwardType = "local"
	// RemoteForward is remote forwarding
	RemoteForward ProxyForwardType = "remote"
	// DynamicForward is dynamic forwarding
	DynamicForward ProxyForwardType = "dynamic"
)

// ProxyForwardConfig is forward: config in proxy;
type ProxyForwardConfig struct {
	Type          ProxyForwardType `json:"type"`
	BindAddress   string           `json:"bindAddress,omitempty"`
	BindPort      int              `json:"bindPort"`
	RemoteAddress string           `json:"remoteAddress,omitempty"`
	RemotePort    int              `json:"remotePort,omitempty"`
}

// NewConfig returns new Config
func NewConfig(source string) (*Config, error) {
	cfg := &Config{}
	err := yaml.Unmarshal([]byte(source), cfg)

	if err != nil {
		return nil, err
	}

	err = SetDefaultConfigValues(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// SetDefaultConfigValues set default values if omitted
func SetDefaultConfigValues(cfg *Config) error {
	if len(cfg.DotFiles.Dir) == 0 {
		cfg.DotFiles.Dir = defaultDotFilesDir
	}

	for i := range cfg.Proxy {
		if cfg.Proxy[i].Port == 0 {
			cfg.Proxy[i].Port = DefaultProxyPort
		}

		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		if strings.HasPrefix(cfg.Proxy[i].PrivateKey, "~/") {
			cfg.Proxy[i].PrivateKey = filepath.Join(home, cfg.Proxy[i].PrivateKey[2:])
		}
	}

	return nil
}

func Ask(param *string, query string, hide bool) error {
	ui := input.DefaultUI()
	ans, err := ui.Ask(query, &input.Options{
		Default:  *param,
		Required: true,
		Loop:     true,
		Hide:     hide,
	})

	if err != nil {
		return err
	}

	*param = strings.TrimSpace(ans)

	return nil
}
