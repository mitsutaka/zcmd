package zcmd

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/cybozu-go/well"
	log "github.com/sirupsen/logrus"
)

const (
	cmdGit = "/usr/bin/git"
)

// Updater is client for repos update
type Updater struct {
	root         string
	repositories []string
}

// NewUpdater returns Updater with given root directory
func NewUpdater(root string) (*Updater, error) {
	if len(root) == 0 {
		return nil, errors.New("no root directory specified")
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

		log.WithFields(log.Fields{
			"path": gitPath,
		}).Info("found git repository")
		u.repositories = append(u.repositories, gitPath)
		return nil
	})
}

// Update fetches, checkouts and pulls git repositories
func (u *Updater) Update(ctx context.Context, jobs int) error {
	jobChan := make(chan struct{}, jobs)
	for i := 0; i < jobs; i++ {
		jobChan <- struct{}{}
	}
	env := well.NewEnvironment(ctx)

	for _, r := range u.repositories {
		ri := r
		env.Go(func(ctx context.Context) error {
			<-jobChan
			defer func() { jobChan <- struct{}{} }()

			err := clean(ctx, ri)
			if err != nil {
				log.WithFields(log.Fields{
					"path":  ri,
					"error": err,
				}).Error("git clean")
				return err
			}
			log.WithFields(log.Fields{
				"path": ri,
			}).Info("git clean")

			err = checkoutMaster(ctx, ri)
			if err != nil {
				log.WithFields(log.Fields{
					"path":  ri,
					"error": err,
				}).Error("git checkout master")
				return err
			}
			log.WithFields(log.Fields{
				"path": ri,
			}).Info("git checkout master")

			err = pull(ctx, ri)
			if err != nil {
				log.WithFields(log.Fields{
					"path":  ri,
					"error": err,
				}).Error("git pull")
				return err
			}
			log.WithFields(log.Fields{
				"path": ri,
			}).Info("git pull")

			err = status(ctx, ri)
			if err != nil {
				log.WithFields(log.Fields{
					"path":  ri,
					"error": err,
				}).Error("git status")
				return err
			}
			log.WithFields(log.Fields{
				"path": ri,
			}).Info("git status")

			return nil
		})
	}

	env.Stop()
	return env.Wait()
}

func clean(ctx context.Context, path string) error {
	args := []string{"clean", "-xdf"}
	cmd := well.CommandContext(ctx, cmdGit, args...)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func checkoutMaster(ctx context.Context, path string) error {
	args := []string{"checkout", "master"}
	cmd := well.CommandContext(ctx, cmdGit, args...)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func pull(ctx context.Context, path string) error {
	args := []string{"pull"}
	cmd := well.CommandContext(ctx, cmdGit, args...)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func status(ctx context.Context, path string) error {
	args := []string{"status"}
	cmd := well.CommandContext(ctx, cmdGit, args...)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
