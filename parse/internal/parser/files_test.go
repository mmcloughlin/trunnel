package parser

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func GlobTest(t *testing.T, pattern string, valid bool) {
	filenames, err := filepath.Glob(pattern)
	require.NoError(t, err)
	for _, filename := range filenames {
		t.Run(filename, func(t *testing.T) {
			_, err := ParseFile(filename)
			if valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestValidFiles(t *testing.T) {
	GlobTest(t, "testdata/valid/*.trunnel", true)
}

func TestFailingFiles(t *testing.T) {
	GlobTest(t, "testdata/failing/*.trunnel", false)
}

func TestTorFiles(t *testing.T) {
	GlobTest(t, "../../../testdata/tor/*.trunnel", true)
}
