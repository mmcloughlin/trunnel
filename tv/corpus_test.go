package tv

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mmcloughlin/trunnel/internal/test"
)

func TestCorpusVectors(t *testing.T) {
	c := &Corpus{}

	assert.Nil(t, c.Vectors("a"))

	vs := []Vector{NewVector([]byte("hello"))}
	c.AddVectors("a", vs)
	assert.Equal(t, vs, c.Vectors("a"))
}

func TestWriteCorpusReal(t *testing.T) {
	c := &Corpus{}
	c.AddVectors("a", []Vector{NewVector([]byte("hello"))})

	// write corpus to a temp directory
	dir, clean := test.TempDir(t)
	defer clean()
	t.Log(dir)

	err := WriteCorpus(c, dir)
	require.NoError(t, err)

	// confirm the vector file exists
	digest := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824" // SHA256("hello")
	path := filepath.Join(dir, "a", digest)
	b, err := ioutil.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, []byte("hello"), b)
}
