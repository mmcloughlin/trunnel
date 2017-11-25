// Package inspect provides tools for extracting information from the AST.
package inspect

import (
	"errors"

	"github.com/mmcloughlin/trunnel/ast"
	"github.com/mmcloughlin/trunnel/fault"
	"github.com/mmcloughlin/trunnel/internal/intervals"
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

// Contexts builds a name to context mapping for all contexts in the file.
func Contexts(f *ast.File) (map[string]*ast.Context, error) {
	ctxs := map[string]*ast.Context{}
	for _, ctx := range f.Contexts {
		n := ctx.Name
		if _, found := ctxs[n]; found {
			return nil, errors.New("duplicate context name")
		}
		ctxs[n] = ctx
	}
	return ctxs, nil
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

// Resolver maintains indexes of various parts of a trunnel file.
type Resolver struct {
	structs   map[string]*ast.Struct
	contexts  map[string]*ast.Context
	constants map[string]int64
}

// NewResolver builds a resolver from the given file.
func NewResolver(f *ast.File) (*Resolver, error) {
	s, err := Structs(f)
	if err != nil {
		return nil, err
	}

	ctxs, err := Contexts(f)
	if err != nil {
		return nil, err
	}

	c, err := Constants(f)
	if err != nil {
		return nil, err
	}

	return &Resolver{
		structs:   s,
		contexts:  ctxs,
		constants: c,
	}, nil
}

// Struct returns the struct with the given name.
func (r *Resolver) Struct(n string) (*ast.Struct, bool) {
	s, ok := r.structs[n]
	return s, ok
}

// Context returns the context with the given name.
func (r *Resolver) Context(n string) (*ast.Context, bool) {
	ctx, ok := r.contexts[n]
	return ctx, ok
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

// Intervals builds an intervals object from an integer list.
func (r *Resolver) Intervals(l *ast.IntegerList) (*intervals.Set, error) {
	is := make([]intervals.Interval, len(l.Ranges))
	for i, rng := range l.Ranges {
		lo, err := r.Integer(rng.Low)
		if err != nil {
			return nil, err
		}
		if rng.High == nil {
			is[i] = intervals.Single(uint64(lo)) // XXX cast
			continue
		}
		hi, err := r.Integer(rng.High)
		if err != nil {
			return nil, err
		}
		is[i] = intervals.Range(uint64(lo), uint64(hi)) // XXX cast
	}
	if intervals.Overlaps(is) {
		return nil, errors.New("overlapping intervals")
	}
	return intervals.NewSet(is...), nil
}

// IntType looks up the integer type refered to by ref. The local struct is
// required to resolve references to fields within the struct.
func (r *Resolver) IntType(ref *ast.IDRef, local *ast.Struct) (*ast.IntType, error) {
	var fs []*ast.Field
	if ref.Scope == "" {
		fs = structFields(local)
	} else {
		ctx, ok := r.Context(ref.Scope)
		if !ok {
			return nil, errors.New("could not find context")
		}
		fs = ctx.Members
	}

	for _, f := range fs {
		if f.Name != ref.Name {
			continue
		}
		i, ok := f.Type.(*ast.IntType)
		if !ok {
			return nil, errors.New("referenced field does not have integer type")
		}
		return i, nil
	}

	return nil, errors.New("could not resolve reference")
}

// structFields extracts the top-level fields from s.
func structFields(s *ast.Struct) []*ast.Field {
	fs := []*ast.Field{}
	for _, m := range s.Members {
		if f, ok := m.(*ast.Field); ok {
			fs = append(fs, f)
		}
	}
	return fs
}
