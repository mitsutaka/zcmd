package zcmd

// Config zcmd config
type Config struct {
	Nas    NasConfig    `yaml:"nas,omitempty"`
	Backup BackupConfig `yaml:"backup,omitempty"`
	Repos  ReposConfig  `yaml:"repos,omitempty"`
}

// NasConfig is nas: config
type NasConfig struct {
	Pull NasPullConfig `yaml:"pull,omitempty"`
	Push NasPushConfig `yaml:"push,omitempty"`
}

// NasPullConfig is pull: config
type NasPullConfig struct {
	Source       string     `yaml:"source"`
	Destinations []PathInfo `yaml:"destinations"`
}

// NasPushConfig is push: config
type NasPushConfig struct {
	Sources     []PathInfo `yaml:"sources"`
	Destination string     `yaml:"destination"`
}

// PathInfo is path information for synchronize directories
type PathInfo struct {
	Name     string   `yaml:"name"`
	Path     string   `yaml:"path"`
	Excludes []string `yaml:"excludes,omitempty"`
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
