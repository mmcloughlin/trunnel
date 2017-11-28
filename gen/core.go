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

func (c Config) write(name string, b []byte) error {
	fn := filepath.Join(c.Dir, name)
	return ioutil.WriteFile(fn, b, 0640)
}

// Package generates a Go package for the given files.
func Package(cfg Config, fs []*ast.File) error {
	// Marshaller file
	b, err := Marshallers(cfg.Package, fs)
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

	corpusDir := filepath.Join(cfg.Dir, "testdata/corpus")
	if err = tv.WriteCorpus(c, corpusDir); err != nil {
		return err
	}

	// Test file
	b, err = CorpusTests(cfg.Package, c)
	if err != nil {
		return err
	}
	if err = cfg.write("gen-marshallers_test.go", b); err != nil {
		return err
	}

	// Fuzzer
	b, err = Fuzzers(cfg.Package, c)
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
