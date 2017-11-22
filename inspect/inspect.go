// Package inspect provides tools for extracting information from the AST.
package inspect

import (
	"errors"

	"github.com/mmcloughlin/trunnel/ast"
	"github.com/mmcloughlin/trunnel/fault"
)

// Structs builds a name to struct mapping for all structs in the file.
func Structs(f *ast.File) (map[string]*ast.Struct, error) {
	structs := map[string]*ast.Struct{}
	for _, s := range f.Structs {
		n := s.Name
		if _, found := structs[n]; found {
			return nil, errors.New("duplicate struct name")
		}
		structs[n] = s
	}
	return structs, nil
}

// Constants builds a map of constant name to value from the declarations in f.
// Errors on duplicate constant names.
func Constants(f *ast.File) (map[string]int64, error) {
	v := map[string]int64{}
	for _, c := range f.Constants {
		n := c.Name
		if _, found := v[n]; found {
			return nil, errors.New("duplicate constant name")
		}
		v[n] = c.Value
	}
	return v, nil
}

// Resolver will resolve AST integers to actual integer values based on defined
// constants.
type Resolver struct {
	constants map[string]int64
}

// NewResolver builds a resolver for the given constants mapping.
func NewResolver(c map[string]int64) *Resolver {
	return &Resolver{
		constants: c,
	}
}

// Integer resolves i to an integer value.
func (r *Resolver) Integer(i ast.Integer) (int64, error) {
	switch i := i.(type) {
	case *ast.IntegerConstRef:
		v, ok := r.constants[i.Name]
		if !ok {
			return 0, errors.New("constant undefined")
		}
		return v, nil
	case *ast.IntegerLiteral:
		return i.Value, nil
	default:
		return 0, fault.NewUnexpectedType(i)
	}
}
