package parse

import (
	"io"
	"strings"

	"github.com/mmcloughlin/trunnel/ast"
	"github.com/mmcloughlin/trunnel/parse/internal/parser"
)

func ParseFile(filename string) (*ast.File, error) {
	return cast(parser.ParseFile(filename))
}

func ParseReader(filename string, r io.Reader) (*ast.File, error) {
	return cast(parser.ParseReader(filename, r))
}

func ParseString(s string) (*ast.File, error) {
	return ParseReader("", strings.NewReader(s))
}

func cast(i interface{}, err error) (*ast.File, error) {
	if err != nil {
		return nil, err
	}
	return i.(*ast.File), nil
}
