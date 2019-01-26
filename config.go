package zcmd

// Config zcmd config
type Config struct {
	Sync   SyncConfig   `yaml:"sync,omitempty"`
	Backup BackupConfig `yaml:"backup,omitempty"`
	Repos  ReposConfig  `yaml:"repos,omitempty"`
}

// SyncConfig is sync: config
type SyncConfig struct {
	Pull []*SyncInfo `yaml:"pull,omitempty"`
	Push []*SyncInfo `yaml:"push,omitempty"`
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

// NewConfig returns new Config
func NewConfig() *Config {
	return &Config{}
}
