package tv

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmcloughlin/trunnel/fault"
	"github.com/mmcloughlin/trunnel/inspect"
	"github.com/mmcloughlin/trunnel/internal/test"
	"github.com/mmcloughlin/trunnel/parse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFiles(t *testing.T) {
	dirs := []string{
		"../testdata/tor",
		"../testdata/trunnel",
	}
	for _, dir := range dirs {
		t.Run(filepath.Base(dir), func(t *testing.T) {
			groups, err := test.LoadFileGroups(dir)
			require.NoError(t, err)
			for _, group := range groups {
				t.Run(strings.Join(group, ","), func(t *testing.T) {
					VerifyGroup(t, group)
				})
			}
		})
	}
}

func VerifyGroup(t *testing.T, filenames []string) {
	fs, err := parse.Files(filenames)
	require.NoError(t, err)

	vs, err := GenerateFiles(fs, WithSelector(RandomSampleSelector(16)))
	if err == fault.ErrNotImplemented {
		t.Log(err)
		t.SkipNow()
	}
	require.NoError(t, err)

	r, err := inspect.NewResolverFiles(fs)
	require.NoError(t, err)

	for _, s := range r.Structs() {
		if s.Extern() {
			continue
		}
		t.Run(s.Name, func(t *testing.T) {
			require.Contains(t, vs, s.Name)
			num := len(vs[s.Name])
			t.Logf("%d test vectors for %s", num, s.Name)
			assert.True(t, num > 0)
		})
	}
}
