// Package parse implements a parser for trunnel source files.
package parse

import (
	"io"
	"strings"

	"github.com/mmcloughlin/trunnel/ast"
	"github.com/mmcloughlin/trunnel/parse/internal/parser"
)

//go:generate pigeon -o internal/parser/gen-parser.go trunnel.pigeon

// File parses filename.
func File(filename string) (*ast.File, error) {
	return cast(parser.ParseFile(filename))
}

// Files is a convenience for parsing multiple files.
func Files(filenames []string) ([]*ast.File, error) {
	fs := make([]*ast.File, len(filenames))
	for i, filename := range filenames {
		f, err := File(filename)
		if err != nil {
			return nil, err
		}
		fs[i] = f
	}
	return fs, nil
}

// Reader parses the data from r using filename as information in
// error messages.
func Reader(filename string, r io.Reader) (*ast.File, error) {
	return cast(parser.ParseReader(filename, r))
}

// String parses s.
func String(s string) (*ast.File, error) {
	return Reader("string", strings.NewReader(s))
}

func cast(i interface{}, err error) (*ast.File, error) {
	if err != nil {
		return nil, err
	}
	return i.(*ast.File), nil
}
