package gen

import (
	"path/filepath"
	"testing"

	"github.com/mmcloughlin/trunnel/parse"
	"github.com/mmcloughlin/trunnel/test"
	"github.com/stretchr/testify/require"
)

func TestFilesBuild(t *testing.T) {
	dirs := []string{
		"testdata/valid",
		"../testdata/tor",
		"../testdata/trunnel",
	}
	for _, d := range dirs {
		pattern := filepath.Join(d, "*.trunnel")
		test.Glob(t, pattern, build)
	}
}

func build(t *testing.T, filename string) {
	f, err := parse.File(filename)
	require.NoError(t, err)

	src, err := File("pkg", f)
	require.NoError(t, err)

	output, err := test.Build(src)
	t.Log(string(output))
	require.NoError(t, err)
}
