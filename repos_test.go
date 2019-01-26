package zcmd

import (
	"context"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	git "gopkg.in/src-d/go-git.v4"
)

//nolint[gochecknoglobals]
var (
	gitURLs = []string{
		"https://github.com/mitsutaka/docker-libs.git",
		"https://github.com/mitsutaka/kubernetes-coreos.git",
		"https://github.com/mitsutaka/submissions.git",
	}
)

func gitClone(dir string) ([]repoInfo, error) {
	repoInfos := make([]repoInfo, len(gitURLs))

	for i, gitURL := range gitURLs {
		u, err := url.Parse(gitURL)
		if err != nil {
			return nil, err
		}
		gitPath := filepath.Join(dir, u.Path)
		gitRepo, err := git.PlainClone(gitPath, false, &git.CloneOptions{
			URL:      gitURL,
			Progress: os.Stdout,
		})
		if err != nil {
			return nil, err
		}
		repoInfos[i].repo = gitRepo
		repoInfos[i].path = gitPath
	}
	return repoInfos, nil
}

func testFind(t *testing.T) {
	t.Parallel()

	// Prepare test data
	dir, err := ioutil.TempDir("", "find")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	repoInfos, err := gitClone(dir)
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
	for _, path := range upd.repositories {
		updPaths = append(updPaths, path.path)
	}
	clonedPaths := make([]string, len(gitURLs))
	for _, path := range repoInfos {
		clonedPaths = append(clonedPaths, path.path)
	}

	if !reflect.DeepEqual(updPaths, clonedPaths) {
		t.Errorf("%v != %v", updPaths, clonedPaths)
	}
}

func testFetch(t *testing.T) {
	t.Parallel()

	// Prepare test data
	dir, err := ioutil.TempDir("", "fetch")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	repoInfos, err := gitClone(dir)
	if err != nil {
		t.Fatal(err)
	}

	for _, repoInfo := range repoInfos {
		err = fetch(context.Background(), repoInfo.repo)
		if err != nil && err != git.NoErrAlreadyUpToDate {
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

	repoInfos, err := gitClone(dir)
	if err != nil {
		t.Fatal(err)
	}

	for _, repoInfo := range repoInfos {
		err = checkout(repoInfo.repo)
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

	repoInfos, err := gitClone(dir)
	if err != nil {
		t.Fatal(err)
	}

	for _, repoInfo := range repoInfos {
		err = pull(context.Background(), repoInfo.repo)
		if err != nil && err != git.NoErrAlreadyUpToDate {
			t.Error(err)
		}
	}
}

func TestReposUpdate(t *testing.T) {
	t.Run("find", testFind)
	t.Run("fetch", testFetch)
	t.Run("checkout", testCheckout)
	t.Run("pull", testPull)
}
