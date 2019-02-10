package zcmd

import (
	yaml "gopkg.in/yaml.v2"
)

// Config zcmd config
type Config struct {
	Sync     SyncConfig     `yaml:"sync,omitempty"`
	Backup   BackupConfig   `yaml:"backup,omitempty"`
	Repos    ReposConfig    `yaml:"repos,omitempty"`
	Dotfiles DotfilesConfig `yaml:"dotfiles,omitempty"`
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
}

// BackupConfig is backup: config
type BackupConfig struct {
	Destinations []string `yaml:"destinations"`
	Includes     []string `yaml:"includes"`
	Excludes     []string `yaml:"excludes,omitempty"`
}

// ReposConfig is repos: config
type ReposConfig struct {
	Root string `yaml:"root"`
}

// DotfilesConfig is dotfiles: config
type DotfilesConfig struct {
	Hosts []string `yaml:"hosts"`
	Files []string `yaml:"files"`
}

// ProxyConfig is proxy: config
type ProxyConfig struct {
	Name       string               `yaml:"name"`
	User       string               `yaml:"user,omitempty"`
	Address    string               `yaml:"address"`
	Port       int                  `yaml:"port,omitempty"`
	PrivateKey string               `yaml:"privateKey"`
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
	BindAddress   string           `yaml:"bindAddress,omitempty"`
	BindPort      int              `yaml:"bindPort"`
	RemoteAddress string           `yaml:"remoteAddress,omitempty"`
	RemotePort    int              `yaml:"remotePort,omitempty"`
}

// NewConfig returns new Config
func NewConfig(source string) (*Config, error) {
	cfg := &Config{}
	err := yaml.Unmarshal([]byte(source), cfg)
	if err != nil {
		return nil, err
	}
	SetDefaultConfigValues(cfg)
	return cfg, nil
}

// SetDefaultConfigValues set default values if omitted
func SetDefaultConfigValues(cfg *Config) {
	for i := range cfg.Proxy {
		if cfg.Proxy[i].Port == 0 {
			cfg.Proxy[i].Port = DefaultProxyPort
		}
	}
}
