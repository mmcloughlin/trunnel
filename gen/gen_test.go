package gen

import (
	"flag"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func (t TestCase) PackageName() string {
	return t.Name
}

func (t TestCase) GoFilename() string {
	return filepath.Join(t.Dir, "gen-"+t.Name+".go")
}

func TestGeneratedFiles(t *testing.T) {
	cases, err := LoadTestCasesGlob("tests/*/*.trunnel")
	require.NoError(t, err)

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			f, err := parse.File(c.TrunnelFile)
			require.NoError(t, err)

			src, err := File(c.PackageName(), f)
			require.NoError(t, err)

			if *update {
				err = ioutil.WriteFile(c.GoFilename(), src, 0640)
				require.NoError(t, err)
			}

			expect, err := ioutil.ReadFile(c.GoFilename())
			require.NoError(t, err)

			assert.Equal(t, expect, src)
		})
	}
}
