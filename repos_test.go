package zcmd

import (
	"context"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/cybozu-go/well"
)

// nolint[gochecknoglobals]
var (
	gitURLs = []string{
		"https://github.com/mitsutaka/docker-libs",
		"https://github.com/mitsutaka/kubernetes-coreos",
	}
)

func gitClone(dir string) ([]string, error) {
	repos := make([]string, len(gitURLs))

	for i, gitURL := range gitURLs {
		u, err := url.Parse(gitURL)
		if err != nil {
			return nil, err
		}

		gitPath := filepath.Join(dir, filepath.Base(u.Path))
		args := []string{"clone", gitURL}
		cmd := well.CommandContext(context.Background(), cmdGit, args...)
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return nil, err
		}

		repos[i] = gitPath
	}

	return repos, nil
}

func testFind(t *testing.T) {
	t.Parallel()

	// Prepare test data
	dir, err := ioutil.TempDir("", "find")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(dir)

	repos, err := gitClone(dir)
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

	updPaths := make([]string, len(gitURLs))

	copy(updPaths, upd.repositories)

	clonedPaths := make([]string, len(gitURLs))

	copy(clonedPaths, repos)

	if !reflect.DeepEqual(updPaths, clonedPaths) {
		t.Errorf("%#v != %#v", updPaths, clonedPaths)
	}
}

func testClean(t *testing.T) {
	t.Parallel()

	// Prepare test data
	dir, err := ioutil.TempDir("", "fetch")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(dir)

	repos, err := gitClone(dir)
	if err != nil {
		t.Fatal(err)
	}

	for _, repo := range repos {
		err = clean(context.Background(), repo)
		if err != nil {
			t.Error(err)
		}
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

	repos, err := gitClone(dir)
	if err != nil {
		t.Fatal(err)
	}

	for _, repo := range repos {
		err = checkoutMaster(context.Background(), repo)
		if err != nil {
			t.Error(err)
		}
	}
}

func testPull(t *testing.T) {
	t.Parallel()

	// Prepare test data
	dir, err := ioutil.TempDir("", "pull")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(dir)

	repos, err := gitClone(dir)
	if err != nil {
		t.Fatal(err)
	}

	for _, repo := range repos {
		err = pull(context.Background(), repo)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestReposUpdate(t *testing.T) {
	t.Run("find", testFind)
	t.Run("clean", testClean)
	t.Run("checkout", testCheckout)
	t.Run("pull", testPull)
}
