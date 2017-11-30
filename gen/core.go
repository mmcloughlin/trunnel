package gen

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/mmcloughlin/random"
	"github.com/serenize/snaker"

	"github.com/mmcloughlin/trunnel/ast"
	"github.com/mmcloughlin/trunnel/fault"
	"github.com/mmcloughlin/trunnel/tv"
)

// Config defines options for package generation.
type Config struct {
	Package string // package name
	Dir     string // directory to output to
	Seed    int64  // seed for corpus generation
}

// PackageName returns the name of the generated package.
func (c Config) PackageName() string {
	return c.Package
}

// OutputDirectory returns the configured output directory.
func (c Config) OutputDirectory() string {
	return c.Dir
}

// Path returns a path to rel within the configured output directory.
func (c Config) Path(rel string) string {
	return filepath.Join(c.OutputDirectory(), rel)
}

func (c Config) write(name string, b []byte) error {
	return ioutil.WriteFile(c.Path(name), b, 0640)
}

// Package generates a Go package for the given files.
func Package(cfg Config, fs []*ast.File) error {
	// Marshaller file
	b, err := Marshallers(cfg.PackageName(), fs)
	if err != nil {
		return err
	}
	if err = cfg.write("gen-marshallers.go", b); err != nil {
		return err
	}

	// Test vector corpus (some features not implemented yet)
	c, err := tv.GenerateFiles(fs, tv.WithRandom(random.NewWithSeed(cfg.Seed)))
	if err == fault.ErrNotImplemented {
		return nil
	}
	if err != nil {
		return err
	}

	corpusDir := filepath.Join(cfg.OutputDirectory(), "testdata/corpus")
	if err = tv.WriteCorpus(c, corpusDir); err != nil {
		return err
	}

	// Test file
	b, err = CorpusTests(cfg.PackageName(), c)
	if err != nil {
		return err
	}
	if err = cfg.write("gen-marshallers_test.go", b); err != nil {
		return err
	}

	// Fuzzer
	b, err = Fuzzers(cfg.PackageName(), c)
	if err != nil {
		return err
	}
	if err = cfg.write("gen-fuzz.go", b); err != nil {
		return err
	}

	return nil
}

func name(n string) string {
	return snaker.SnakeToCamel(strings.ToLower(n))
}
