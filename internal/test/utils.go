// Package test provides utilities for trunnel testing.
package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

// Glob runs a test on all the files matching a glob pattern.
func Glob(t *testing.T, pattern string, f func(*testing.T, string)) {
	filenames, err := filepath.Glob(pattern)
	require.NoError(t, err)
	for _, filename := range filenames {
		t.Run(filename, func(t *testing.T) {
			f(t, filename)
		})
	}
}

// TrunnelFiles returns all the trunnel files in the given directory.
func TrunnelFiles(dir string) ([]string, error) {
	pattern := filepath.Join(dir, "*.trunnel")
	return filepath.Glob(pattern)
}

// LoadFileGroups looks for trunnel files in a directory and returns groups of
// files that can be "compiled" together (accounting for extern struct
// declarations). Dependencies can be recorded in a deps.yaml file in the
// directory.
func LoadFileGroups(dir string) ([][]string, error) {
	deps, err := LoadDependenciesDir(dir)
	if err != nil {
		return nil, err
	}

	filenames, err := TrunnelFiles(dir)
	if err != nil {
		return nil, err
	}

	groups := [][]string{}
	for _, filename := range filenames {
		group := []string{filename}
		base := filepath.Base(filename)
		if ds, ok := deps.Dependencies[base]; ok {
			for _, d := range ds {
				group = append(group, filepath.Join(dir, d))
			}
		}
		groups = append(groups, group)
	}

	return groups, nil
}

// Dependencies records dependencies between trunnel files.
type Dependencies struct {
	Dependencies map[string][]string
}

// LoadDependenciesFile loads Dependencies from a YAML file.
func LoadDependenciesFile(filename string) (*Dependencies, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	deps := &Dependencies{}
	if err := yaml.Unmarshal(b, deps); err != nil {
		return nil, err
	}
	return deps, nil
}

// LoadDependenciesDir looks for "deps.yml" in the directory and loads it if
// it exists. If the file is not found, it loads an empty set of dependencies.
func LoadDependenciesDir(dir string) (*Dependencies, error) {
	filename := filepath.Join(dir, "deps.yml")
	deps, err := LoadDependenciesFile(filename)
	if os.IsNotExist(err) {
		return &Dependencies{
			Dependencies: map[string][]string{},
		}, nil
	}
	return deps, err
}

// TempDir creates a temp directory. Returns the path to the directory and a
// cleanup function.
func TempDir(t *testing.T) (string, func()) {
	dir, err := ioutil.TempDir("", "trunnel")
	require.NoError(t, err)
	return dir, func() {
		require.NoError(t, os.RemoveAll(dir))
	}
}

// Build checks whether Go source code src builds correctly. Returns the output
// of "go build" and an error, if any.
func Build(srcs [][]byte) ([]byte, error) {
	dir, err := ioutil.TempDir("", "trunnel")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	filenames := []string{}
	for i, src := range srcs {
		filename := filepath.Join(dir, fmt.Sprintf("src%03d.go", i))
		if err := ioutil.WriteFile(filename, src, 0600); err != nil {
			return nil, err
		}
		filenames = append(filenames, filename)
	}

	args := append([]string{"build"}, filenames...)
	cmd := exec.Command("go", args...)
	return cmd.CombinedOutput()
}

// FileContentsEqual determines whether a and b have the same contents.
func FileContentsEqual(a, b string) (bool, error) {
	da, err := ioutil.ReadFile(a)
	if err != nil {
		return false, err
	}

	db, err := ioutil.ReadFile(b)
	if err != nil {
		return false, err
	}

	return bytes.Equal(da, db), nil
}

// AssertFileContentsEqual asserts that files a and b have the same contents.
func AssertFileContentsEqual(t *testing.T, a, b string) {
	eq, err := FileContentsEqual(a, b)
	require.NoError(t, err)
	assert.True(t, eq)
}
