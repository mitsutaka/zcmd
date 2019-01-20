package zcmd

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/cybozu-go/well"
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

		log.Printf("found %s", gitPath)
		repoInfo := repoInfo{
			path: gitPath,
			repo: repo,
		}
		u.repositories = append(u.repositories, repoInfo)
		return nil
	})
}

// FetchRepositories fetches git repositories
func (u *Updater) FetchRepositories(ctx context.Context) {
	env := well.NewEnvironment(ctx)
	for _, r := range u.repositories {
		ri := r
		env.Go(func(ctx context.Context) error {
			err := fetch(ctx, ri.repo)
			if err != nil && err != git.NoErrAlreadyUpToDate {
				log.Printf("git fetch %s, error: %#v\n", ri.path, err)
				return err
			}

			log.Printf("git fetched %s\n", ri.path)
			return nil
		})
	}
	env.Stop()
	_ = env.Wait()
}

func fetch(ctx context.Context, repo *git.Repository) error {
	remote, err := repo.Remote("origin")
	if err != nil {
		return err
	}

	return remote.FetchContext(ctx, &git.FetchOptions{
		Progress: os.Stdout,
	})
}

// CheckoutRepositories checkouts latest commits
func (u *Updater) CheckoutRepositories(ctx context.Context) {
	env := well.NewEnvironment(ctx)
	for _, r := range u.repositories {
		ri := r
		env.Go(func(ctx context.Context) error {
			err := checkout(ri.repo)
			if err != nil {
				log.Printf("git fetch %s, error: %#v\n", ri.path, err)
				return err
			}
			log.Printf("git checked out %s\n", ri.path)
			return nil
		})
	}
	env.Stop()
	_ = env.Wait()
}

func checkout(repo *git.Repository) error {
	head, err := repo.Head()
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	return worktree.Checkout(&git.CheckoutOptions{
		Hash: head.Hash(),
	})
}
