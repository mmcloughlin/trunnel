package tv

import (
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// Corpus is a collection of test vectors for multiple types.
type Corpus struct {
	Suites []Suite
}

// Suite contains test vectors for one type.
type Suite struct {
	Type    string
	Vectors []Vector
}

// AddSuite appends a suite to the corpus.
func (c *Corpus) AddSuite(s Suite) {
	c.Suites = append(c.Suites, s)
}

// AddVectors adds vectors to the corpus for type n. Convenience wrapper around
// AddSuite.
func (c *Corpus) AddVectors(n string, vs []Vector) {
	c.AddSuite(Suite{
		Type:    n,
		Vectors: vs,
	})
}

// Vectors looks up test vectors for the given type. Returns nil if none found.
func (c Corpus) Vectors(n string) []Vector {
	for _, s := range c.Suites {
		if s.Type == n {
			return s.Vectors
		}
	}
	return nil
}

// WriteCorpus writes the corpus of vectors in a standard structure under dir.
func WriteCorpus(c *Corpus, dir string) error {
	fs := afero.NewBasePathFs(afero.NewOsFs(), dir)
	return writecorpus(c, fs, sha256namer)
}

// corpus writes vectors to the filesystem fs.
func writecorpus(c *Corpus, fs afero.Fs, namer func([]byte) string) error {
	a := afero.Afero{Fs: fs}
	for _, s := range c.Suites {
		dir := s.Type
		if err := a.MkdirAll(dir, 0750); err != nil {
			return errors.Wrap(err, "could not create directory")
		}

		for _, v := range s.Vectors {
			filename := namer(v.Data)
			path := filepath.Join(dir, filename)
			if err := a.WriteFile(path, v.Data, 0640); err != nil {
				return errors.Wrap(err, "failed to write file")
			}
		}
	}
	return nil
}

// sha256namer returns the hex-encoded sha256 hash of b.
func sha256namer(b []byte) string {
	d := sha256.Sum256(b)
	return hex.EncodeToString(d[:])
}
