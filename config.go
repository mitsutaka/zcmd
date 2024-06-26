package zcmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/tcnksm/go-input"
)

type Config struct {
	Sync     SyncConfig     `yaml:"sync,omitempty"`
	Backup   BackupConfig   `yaml:"backup,omitempty"`
	Repos    ReposConfig    `yaml:"repos,omitempty"`
	DotFiles DotFilesConfig `yaml:"dotfiles,omitempty"`
	Proxy    []ProxyConfig  `yaml:"proxy,omitempty"`
}

// SyncConfig is sync: config
type SyncConfig struct {
	Pull []SyncInfo `yaml:"pull,omitempty"`
	Push []SyncInfo `yaml:"push,omitempty"`
}

// SyncInfo is path information for synchronize directories
type SyncInfo struct {
	Name        string   `yaml:"name"`
	Source      string   `yaml:"source"`
	Destination string   `yaml:"destination"`
	Excludes    []string `yaml:"excludes,omitempty"`
	DisableSudo bool     `yaml:"disable_sudo,omitempty"`
}

// BackupConfig is backup: config
type BackupConfig struct {
	Destinations []string `yaml:"destinations"`
	Includes     []string `yaml:"includes"`
	Excludes     []string `yaml:"excludes,omitempty"`
	BackupPrefix string   `yaml:"prefix,omitempty"`
}

// ReposConfig is repos: config
type ReposConfig struct {
	Root string `yaml:"root"`
}

// nolint[gochecknoglobals]
// defaultDotFilesDir is default dotfiles directory path
var defaultDotFilesDir = filepath.Join(os.Getenv("HOME"), ".zdotfiles")

// DotFilesConfig is dotfiles: config
type DotFilesConfig struct {
	Dir   string   `yaml:"dir,omitempty"`
	Hosts []string `yaml:"hosts"`
	Files []string `yaml:"files"`
}

// ProxyConfig is proxy: config
type ProxyConfig struct {
	Name       string               `yaml:"name"`
	User       string               `yaml:"user,omitempty"`
	Address    string               `yaml:"address"`
	Port       int                  `yaml:"port,omitempty"`
	PrivateKey string               `yaml:"private_key"`
	Forward    []ProxyForwardConfig `yaml:"forward"`
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
	Type          ProxyForwardType `yaml:"type"`
	BindAddress   string           `yaml:"bind_address,omitempty"`
	BindPort      int              `yaml:"bind_port"`
	RemoteAddress string           `yaml:"remote_address,omitempty"`
	RemotePort    int              `yaml:"remote_port,omitempty"`
}

// NewConfig returns new Config
func NewConfig(source []byte) (*Config, error) {
	cfg := &Config{}

	viper.SetConfigType("yaml")

	err := viper.ReadConfig(bytes.NewBuffer(source))
	if err != nil {
		return nil, err
	}

	yamlTagOption := func(c *mapstructure.DecoderConfig) {
		c.TagName = "yaml"
	}

	err = viper.Unmarshal(&cfg, yamlTagOption)
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
