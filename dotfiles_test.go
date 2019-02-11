package zcmd

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

//nolint[gochecknoglobals]
const (
	gitURL = "https://github.com/mitsutaka/docker-libs.git"
)

func TestDotfiles(t *testing.T) {
	t.Parallel()

	// Prepare test data
	dir, err := ioutil.TempDir("", "init")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	err = os.Setenv("HOME", filepath.Join(dir, "home"))
	if err != nil {
		t.Fatal(err)
	}

	df, err := NewDotFiler(&DotFilesConfig{
		Dir:   dir,
		Files: []string{"znc/Dockerfile"},
	})
	if err != nil {
		t.Error(err)
	}

	err = df.Init(context.Background(), gitURL)
	if err != nil {
		t.Error(err)
	}

	_, err = os.Stat(filepath.Join(dir, "home", ".znc", "Dockerfile"))
	if err != nil {
		t.Error(err)
	}

	err = os.Remove(filepath.Join(dir, "home", ".znc", "Dockerfile"))
	if err != nil {
		t.Fatal(err)
	}

	err = df.Pull(context.Background())
	if err != nil {
		t.Error(err)
	}

	_, err = os.Stat(filepath.Join(dir, "home", ".znc", "Dockerfile"))
	if err != nil {
		t.Error(err)
	}
}
