package gen

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mmcloughlin/trunnel/internal/test"
	"github.com/mmcloughlin/trunnel/parse"
)

func TestFilesBuild(t *testing.T) {
	dirs := []string{
		"testdata/valid",
		"../testdata/tor",
		"../testdata/trunnel",
	}
	for _, dir := range dirs {
		t.Run(filepath.Base(dir), func(t *testing.T) {
			groups, err := test.LoadFileGroups(dir)
			require.NoError(t, err)
			for _, group := range groups {
				t.Run(strings.Join(group, ","), func(t *testing.T) {
					Build(t, group)
				})
			}
		})
	}
}

func Build(t *testing.T, filenames []string) {
	srcs := [][]byte{}
	for _, filename := range filenames {
		f, err := parse.File(filename)
		require.NoError(t, err)
		src, err := File("pkg", f)
		require.NoError(t, err)
		srcs = append(srcs, src)
	}

	output, err := test.Build(srcs)
	if err != nil {
		t.Fatal(string(output))
	}
}
