package parse

import (
	"github.com/mmcloughlin/trunnel/ast"
	"github.com/mmcloughlin/trunnel/parse/internal/parser"
)

func ParseFile(filename string) (*ast.File, error) {
	i, err := parser.ParseFile(filename)
	if err != nil {
		return nil, err
	}
	return i.(*ast.File), nil
}
