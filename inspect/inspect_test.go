package inspect

import (
	"testing"

	"github.com/mmcloughlin/trunnel/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStructs(t *testing.T) {
	f := &ast.File{
		Structs: []*ast.Struct{
			&ast.Struct{Name: "a"},
			&ast.Struct{Name: "b"},
			&ast.Struct{Name: "c"},
		},
	}
	s, err := Structs(f)
	require.NoError(t, err)
	assert.Equal(t, map[string]*ast.Struct{
		"a": &ast.Struct{Name: "a"},
		"b": &ast.Struct{Name: "b"},
		"c": &ast.Struct{Name: "c"},
	}, s)
}

func TestStructsDupe(t *testing.T) {
	f := &ast.File{
		Structs: []*ast.Struct{
			&ast.Struct{Name: "a"},
			&ast.Struct{Name: "b"},
			&ast.Struct{Name: "a"},
		},
	}
	_, err := Structs(f)
	assert.EqualError(t, err, "duplicate struct name")
}

func TestConstants(t *testing.T) {
	f := &ast.File{
		Constants: []*ast.Constant{
			&ast.Constant{Name: "a", Value: 1},
			&ast.Constant{Name: "b", Value: 2},
			&ast.Constant{Name: "c", Value: 3},
		},
	}
	v, err := Constants(f)
	require.NoError(t, err)
	assert.Equal(t, map[string]int64{
		"a": 1,
		"b": 2,
		"c": 3,
	}, v)
}

func TestConstantsDupe(t *testing.T) {
	f := &ast.File{
		Constants: []*ast.Constant{
			&ast.Constant{Name: "a"},
			&ast.Constant{Name: "b"},
			&ast.Constant{Name: "a"},
		},
	}
	_, err := Constants(f)
	assert.EqualError(t, err, "duplicate constant name")
}

func TestNewResolverErrors(t *testing.T) {
	files := []*ast.File{
		&ast.File{
			Structs: []*ast.Struct{
				&ast.Struct{Name: "a"},
				&ast.Struct{Name: "a"},
			},
		},
		&ast.File{
			Contexts: []*ast.Context{
				&ast.Context{Name: "a"},
				&ast.Context{Name: "a"},
			},
		},
		&ast.File{
			Constants: []*ast.Constant{
				&ast.Constant{Name: "a"},
				&ast.Constant{Name: "a"},
			},
		},
	}
	for _, f := range files {
		_, err := NewResolver(f)
		assert.Error(t, err)
	}
}

func TestResolverStruct(t *testing.T) {
	f := &ast.File{
		Structs: []*ast.Struct{
			&ast.Struct{Name: "a"},
			&ast.Struct{Name: "b"},
		},
	}
	r, err := NewResolver(f)
	require.NoError(t, err)

	s, ok := r.Struct("b")
	assert.True(t, ok)
	assert.Equal(t, &ast.Struct{Name: "b"}, s)

	_, ok = r.Struct("idk")
	assert.False(t, ok)
}

func TestResolverContext(t *testing.T) {
	f := &ast.File{
		Contexts: []*ast.Context{
			&ast.Context{Name: "a"},
			&ast.Context{Name: "b"},
		},
	}
	r, err := NewResolver(f)
	require.NoError(t, err)

	c, ok := r.Context("b")
	assert.True(t, ok)
	assert.Equal(t, &ast.Context{Name: "b"}, c)

	_, ok = r.Context("idk")
	assert.False(t, ok)
}

func TestResolverInteger(t *testing.T) {
	f := &ast.File{
		Constants: []*ast.Constant{
			&ast.Constant{Name: "a", Value: 1},
			&ast.Constant{Name: "b", Value: 2},
			&ast.Constant{Name: "c", Value: 3},
		},
	}
	r, err := NewResolver(f)
	require.NoError(t, err)

	cases := []struct {
		Name     string
		Int      ast.Integer
		Value    int64
		HasError bool
	}{
		{
			Name:  "constref",
			Int:   &ast.IntegerConstRef{Name: "b"},
			Value: 2,
		},
		{
			Name:     "undefconst",
			Int:      &ast.IntegerConstRef{Name: "no"},
			HasError: true,
		},
		{
			Name:  "literal",
			Int:   &ast.IntegerLiteral{Value: 42},
			Value: 42,
		},
		{
			Name:     "unexpectedtype",
			Int:      &ast.File{},
			HasError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			v, err := r.Integer(c.Int)
			assert.Equal(t, c.Value, v)
			assert.Equal(t, c.HasError, err != nil)
		})
	}
}

func TestResolverIntType(t *testing.T) {
	s := &ast.Struct{
		Name: "name",
		Members: []ast.Member{
			&ast.Field{Name: "a", Type: ast.U8},
			&ast.Field{Name: "b", Type: ast.U16},
			&ast.Field{Name: "c", Type: ast.U32},
			&ast.Field{Name: "s", Type: &ast.NulTermString{}},
		},
	}
	f := &ast.File{
		Structs: []*ast.Struct{s},
		Contexts: []*ast.Context{
			&ast.Context{
				Name: "ctx",
				Members: []*ast.Field{
					&ast.Field{Name: "a", Type: ast.U8},
					&ast.Field{Name: "b", Type: ast.U16},
					&ast.Field{Name: "c", Type: ast.U32},
				},
			},
		},
	}
	r, err := NewResolver(f)
	require.NoError(t, err)

	cases := []struct {
		Name     string
		Ref      *ast.IDRef
		IntType  *ast.IntType
		HasError bool
	}{
		{
			Name:    "local",
			Ref:     &ast.IDRef{Name: "b"},
			IntType: ast.U16,
		},
		{
			Name:    "ctx",
			Ref:     &ast.IDRef{Scope: "ctx", Name: "c"},
			IntType: ast.U32,
		},
		{
			Name:     "undefctx",
			Ref:      &ast.IDRef{Scope: "what", Name: "c"},
			HasError: true,
		},
		{
			Name:     "undeffield",
			Ref:      &ast.IDRef{Scope: "ctx", Name: "missing"},
			HasError: true,
		},
		{
			Name:     "notint",
			Ref:      &ast.IDRef{Name: "s"},
			HasError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			i, err := r.IntType(c.Ref, s)
			assert.Equal(t, c.IntType, i)
			assert.Equal(t, c.HasError, err != nil)
		})
	}
}
