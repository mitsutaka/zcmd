package zcmd

import (
	"context"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"
)

// DotFiler is client for dotfiles management
type DotFiler struct {
	homeDir  string
	hostname string
	config   *DotFilesConfig
}

// NewDotFiler returns DotFiler with given home directory
func NewDotFiler(cfg *DotFilesConfig) (*DotFiler, error) {
	_, err := os.Stat(cfg.Dir)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &DotFiler{
		homeDir:  os.Getenv("HOME"),
		hostname: hostname,
		config:   cfg,
	}, nil
}

// Init initializes dotfiles directory if not exist
func (d *DotFiler) Init(ctx context.Context, gitURL string) error {
	_, err := os.Stat(d.config.Dir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	_, err = git.PlainCloneContext(ctx, d.config.Dir, false, &git.CloneOptions{
		URL:      gitURL,
		Progress: os.Stdout,
	})
	if err != nil {
		return err
	}

	return d.MakeSymlinks()
}

// Pull updates dotfiles and make symlinks
func (d *DotFiler) Pull(ctx context.Context) error {
	repo, err := git.PlainOpen(d.config.Dir)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"path": d.config.Dir,
	}).Info("found git repository")

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = worktree.PullContext(ctx, &git.PullOptions{
		Progress: os.Stdout,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	return d.MakeSymlinks()
}

// MakeSymlinks makes symlinks to the home directory
func (d *DotFiler) MakeSymlinks() error {
	for _, p := range d.config.Files {
		src := filepath.Join(d.config.Dir, p)
		dst := filepath.Join(d.homeDir, "."+p)
		dstDir := filepath.Dir(dst)

		err := os.MkdirAll(dstDir, 0755)
		if err != nil {
			return err
		}

		_, err = os.Stat(dst)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		if err == nil {
			err := os.Remove(dst)
			if err != nil {
				return err
			}
		}

		err = os.Symlink(src, dst)
		if err != nil {
			return err
		}
		log.WithFields(log.Fields{
			"src": src,
			"dst": dst,
		}).Info("symlinked")
	}
	return nil
}
