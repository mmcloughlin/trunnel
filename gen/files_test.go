package gen

import (
	"path/filepath"
	"testing"

	"github.com/mmcloughlin/trunnel/internal/test"
	"github.com/mmcloughlin/trunnel/parse"
	"github.com/stretchr/testify/require"
)

func TestFilesBuild(t *testing.T) {
	cases := []struct {
		Name         string
		Dir          string
		Dependencies map[string][]string
	}{
		{
			Name: "valid",
			Dir:  "testdata/valid",
		},
		{
			Name: "tor",
			Dir:  "../testdata/tor",
			Dependencies: map[string][]string{
				"cell_establish_intro.trunnel": []string{"cell_common.trunnel"},
				"cell_introduce1.trunnel":      []string{"cell_common.trunnel", "ed25519_cert.trunnel"},
			},
		},
		{
			Name: "trunnel",
			Dir:  "../testdata/trunnel",
			Dependencies: map[string][]string{
				"derived.trunnel": []string{"simple.trunnel"},
				"opaque.trunnel":  []string{"simple.trunnel"},
				"prop224.trunnel": []string{"prop220.trunnel"},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			pattern := filepath.Join(c.Dir, "*.trunnel")
			test.Glob(t, pattern, func(t *testing.T, filename string) {
				filenames := []string{filename}
				base := filepath.Base(filename)
				if deps, ok := c.Dependencies[base]; ok {
					for _, dep := range deps {
						filenames = append(filenames, filepath.Join(c.Dir, dep))
					}
				}
				Build(t, filenames)
			})
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
