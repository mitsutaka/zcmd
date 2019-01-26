package zcmd

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/cybozu-go/well"
	log "github.com/sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"
)

// Updater is client for repos update
type Updater struct {
	root         string
	repositories []repoInfo
}

type repoInfo struct {
	path string
	repo *git.Repository
}

// NewUpdater returns Updater with given root directory
func NewUpdater(root string) (*Updater, error) {
	if len(root) == 0 {
		return nil, errors.New("Not root directory specified")
	}
	return &Updater{
		root: root,
	}, nil
}

// FindRepositories traverses given directory and return git repositories
func (u *Updater) FindRepositories() error {
	return filepath.Walk(u.root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return err
		}
		if info.Name() != ".git" {
			// Skip
			return nil
		}
		gitPath := filepath.Dir(path)

		repo, err := git.PlainOpen(gitPath)
		if err != nil {
			// Skip
			return err
		}

		log.WithFields(log.Fields{
			"path": gitPath,
		}).Info("found git repository")
		repoInfo := repoInfo{
			path: gitPath,
			repo: repo,
		}
		u.repositories = append(u.repositories, repoInfo)
		return nil
	})
}

// Update fetches, checkouts and pulls git repositories
func (u *Updater) Update(ctx context.Context) error {
	env := well.NewEnvironment(ctx)

	for _, r := range u.repositories {
		ri := r
		env.Go(func(ctx context.Context) error {
			err := fetch(ctx, ri.repo)
			if err != nil && err != git.NoErrAlreadyUpToDate {
				log.WithFields(log.Fields{
					"command": "git fetch",
					"path":    ri.path,
					"error":   err,
				}).Error("fetched")
				return err
			}
			log.WithFields(log.Fields{
				"command": "git fetch",
				"path":    ri.path,
				"error":   err,
			}).Info("fetched")

			err = checkout(ri.repo)
			if err != nil {
				log.WithFields(log.Fields{
					"command": "git checkout",
					"path":    ri.path,
					"error":   err,
				}).Error("checked out")
				return err
			}
			log.WithFields(log.Fields{
				"command": "git checkout",
				"path":    ri.path,
				"error":   err,
			}).Info("checked out")

			err = pull(ctx, ri.repo)
			if err != nil && err != git.NoErrAlreadyUpToDate {
				log.WithFields(log.Fields{
					"command": "git pull",
					"path":    ri.path,
					"error":   err,
				}).Error("pulled")
				return err
			}
			log.WithFields(log.Fields{
				"command": "git pull",
				"path":    ri.path,
				"error":   err,
			}).Info("pulled")

			return nil
		})
	}

	env.Stop()
	return env.Wait()
}

func fetch(ctx context.Context, repo *git.Repository) error {
	remote, err := repo.Remote("origin")
	if err != nil {
		return err
	}

	return remote.FetchContext(ctx, &git.FetchOptions{})
}

func checkout(repo *git.Repository) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	// Checkout master branch
	return worktree.Checkout(&git.CheckoutOptions{})
}

func pull(ctx context.Context, repo *git.Repository) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	return worktree.PullContext(ctx, &git.PullOptions{
		Progress: os.Stdout,
	})
}
