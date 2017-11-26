package parse

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileNotExists(t *testing.T) {
	_, err := File("doesnotexist")
	assert.Error(t, err)
}

func TestFilesNotExists(t *testing.T) {
	_, err := Files([]string{"doesnotexist"})
	assert.Error(t, err)
}

func TestFiles(t *testing.T) {
	filenames, err := filepath.Glob("testdata/valid/*.trunnel")
	require.NoError(t, err)
	fs, err := Files(filenames)
	require.NoError(t, err)
	assert.Len(t, fs, 3)
}

func TestErrorReader(t *testing.T) {
	_, err := Reader("", errorReader{})
	assert.Error(t, err)
}

type errorReader struct{}

func (r errorReader) Read(_ []byte) (int, error) {
	return 0, assert.AnError
}
