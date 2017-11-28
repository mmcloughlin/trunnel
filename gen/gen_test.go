package gen

import (
	"flag"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mmcloughlin/trunnel/ast"
	"github.com/mmcloughlin/trunnel/internal/test"
	"github.com/mmcloughlin/trunnel/parse"
)

var update = flag.Bool("update", false, "update golden files")

type TestCase struct {
	TrunnelFile string
	Dir         string
	Name        string
}

func NewTestCaseFromTrunnel(path string) TestCase {
	dir, file := filepath.Split(path)
	name := strings.TrimSuffix(file, filepath.Ext(file))
	return TestCase{
		TrunnelFile: path,
		Dir:         dir,
		Name:        name,
	}
}

func LoadTestCasesGlob(pattern string) ([]TestCase, error) {
	filenames, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	t := make([]TestCase, len(filenames))
	for i, filename := range filenames {
		t[i] = NewTestCaseFromTrunnel(filename)
	}
	return t, nil
}

func (t TestCase) Config() Config {
	return Config{
		Package: t.Name,
		Dir:     t.Dir,
	}
}

func TestGeneratedFiles(t *testing.T) {
	cases, err := LoadTestCasesGlob("tests/*/*.trunnel")
	require.NoError(t, err)

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			tmp, clean := test.TempDir(t)
			defer clean()

			f, err := parse.File(c.TrunnelFile)
			require.NoError(t, err)

			cfg := c.Config()
			if !*update {
				cfg.Dir = tmp
			}
			err = Package(cfg, []*ast.File{f})
			require.NoError(t, err)

			cmp := []string{
				"gen-marshallers.go",
				"gen-marshallers_test.go",
			}
			for _, path := range cmp {
				t.Run(path, func(t *testing.T) {
					got := filepath.Join(cfg.Dir, path)
					expect := filepath.Join(c.Dir, path)
					if !test.FileExists(expect) {
						t.SkipNow()
					}
					test.AssertFileContentsEqual(t, expect, got)
				})
			}
		})
	}
}
