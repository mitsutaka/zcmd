package zcmd

import (
	"context"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	git "gopkg.in/src-d/go-git.v4"
)

var (
	gitURLs = []string{
		"https://github.com/mitsutaka/docker-libs.git",
		"https://github.com/mitsutaka/kubernetes-coreos.git",
		"https://github.com/mitsutaka/submissions.git",
	}
)

func gitClone(dir string) error {
	for _, gitURL := range gitURLs {
		u, err := url.Parse(gitURL)
		if err != nil {
			return err
		}
		_, err = git.PlainClone(filepath.Join(dir, u.Path), false, &git.CloneOptions{
			URL:      gitURL,
			Progress: os.Stdout,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func testFetch(t *testing.T) {
	t.Parallel()

	// Prepare test data
	dir, err := ioutil.TempDir("", "fetch")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	err = gitClone(dir)
	if err != nil {
		t.Fatal(err)
	}

	upd, err := NewUpdater(dir)
	if err != nil {
		t.Error(err)
	}

	err = upd.FindRepositories()
	if err != nil {
		t.Error(err)
	}

	err = upd.FetchRepositories(context.Background())
	if err != nil {
		t.Error(err)
	}
}

func testCheckout(t *testing.T) {
	t.Parallel()

	// Prepare test data
	dir, err := ioutil.TempDir("", "checkout")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	err = gitClone(dir)
	if err != nil {
		t.Fatal(err)
	}

	upd, err := NewUpdater(dir)
	if err != nil {
		t.Error(err)
	}

	err = upd.FindRepositories()
	if err != nil {
		t.Error(err)
	}

	err = upd.CheckoutRepositories(context.Background())
	if err != nil {
		t.Error(err)
	}
}

func TestReposUpdate(t *testing.T) {
	t.Run("fetch", testFetch)
	t.Run("checkout", testCheckout)
}
