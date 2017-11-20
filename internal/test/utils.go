// Package test provides utilities for trunnel testing.
package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
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
