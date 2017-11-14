package gen

import (
	"flag"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmcloughlin/trunnel/parse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update golden files")

func PackageName(path string) string {
	return filepath.Base(filepath.Dir(path))
}

func GoFilename(path string) string {
	return strings.TrimSuffix(path, filepath.Ext(path)) + ".go"
}

func TestGeneratedFiles(t *testing.T) {
	filenames, err := filepath.Glob("tests/*/*.trunnel")
	require.NoError(t, err)
	for _, filename := range filenames {
		t.Run(filename, func(t *testing.T) {
			f, err := parse.File(filename)
			require.NoError(t, err)

			pkg := PackageName(filename)
			src, err := File(pkg, f)
			require.NoError(t, err)

			if *update {
				err = ioutil.WriteFile(GoFilename(filename), src, 0640)
				require.NoError(t, err)
			}

			expect, err := ioutil.ReadFile(GoFilename(filename))
			require.NoError(t, err)

			assert.Equal(t, expect, src)
		})
	}
}
