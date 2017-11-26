package inspect

import (
	"testing"

	"github.com/mmcloughlin/trunnel/ast"
	"github.com/mmcloughlin/trunnel/internal/intervals"
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

func TestResolverAddFileErrors(t *testing.T) {
	files := []*ast.File{
		&ast.File{
			Structs: []*ast.Struct{
				&ast.Struct{Name: "a", Members: []ast.Member{
					&ast.Field{Name: "x", Type: ast.U8},
				}},
			}},
		&ast.File{Contexts: []*ast.Context{&ast.Context{Name: "a"}}},
		&ast.File{Constants: []*ast.Constant{&ast.Constant{Name: "a"}}},
	}
	for _, f := range files {
		r := NewResolverEmpty()
		err := r.AddFile(f)
		require.NoError(t, err)
		f.Constants = []*ast.Constant{{Name: "a", Value: 1}} // hack to force constant conflict
		err = r.AddFile(f)
		require.Error(t, err)
	}
}

func TestResolverAddStructOverrideError(t *testing.T) {
	r := NewResolverEmpty()

	s := &ast.Struct{Name: "a", Members: []ast.Member{
		&ast.Field{Name: "x", Type: ast.U16},
	}}

	err := r.AddStruct(s)
	require.NoError(t, err)

	err = r.AddStruct(s)
	require.EqualError(t, err, "cannot override non-extern struct")
}

func TestResolverAddStructOverrideExtern(t *testing.T) {
	r := NewResolverEmpty()

	err := r.AddStruct(&ast.Struct{Name: "a"})
	require.NoError(t, err)

	err = r.AddStruct(&ast.Struct{Name: "a", Members: []ast.Member{
		&ast.Field{Name: "x", Type: ast.U16},
	}})
	require.NoError(t, err)
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

func TestResolverStructNonExtern(t *testing.T) {
	f := &ast.File{
		Structs: []*ast.Struct{
			&ast.Struct{
				Name: "a",
				Members: []ast.Member{
					&ast.Field{Name: "x", Type: ast.U16},
				},
			},
			&ast.Struct{
				Name:    "b",
				Members: nil, // extern
			},
		},
	}
	r, err := NewResolver(f)
	require.NoError(t, err)

	_, err = r.StructNonExtern("a")
	assert.NoError(t, err)

	_, err = r.StructNonExtern("b")
	assert.EqualError(t, err, "struct is external")

	_, err = r.StructNonExtern("c")
	assert.EqualError(t, err, "struct not found")
}

func TestResolverAddContext(t *testing.T) {
	r := NewResolverEmpty()
	ctx := &ast.Context{Name: "a"}

	err := r.AddContext(ctx)
	require.NoError(t, err)

	err = r.AddContext(ctx)
	require.EqualError(t, err, "cannot override context")

	c, ok := r.Context("a")
	assert.True(t, ok)
	assert.Equal(t, &ast.Context{Name: "a"}, c)
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

func TestResolverAddConstant(t *testing.T) {
	r := NewResolverEmpty()
	c := &ast.Constant{
		Name:  "X",
		Value: 42,
	}
	err := r.AddConstant(c)
	require.NoError(t, err)

	v, ok := r.Constant("X")
	require.True(t, ok)
	assert.Equal(t, int64(42), v)
}

func TestResolverSetConstantTwice(t *testing.T) {
	r := NewResolverEmpty()
	err := r.SetConstant("X", 42)
	require.NoError(t, err)

	err = r.SetConstant("X", 42)
	require.NoError(t, err)

	v, ok := r.Constant("X")
	require.True(t, ok)
	assert.Equal(t, int64(42), v)
}

func TestResolverSetConstantOverride(t *testing.T) {
	r := NewResolverEmpty()
	err := r.SetConstant("X", 42)
	require.NoError(t, err)

	err = r.SetConstant("X", 43)
	assert.Error(t, err)
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

func TestResolverIntervals(t *testing.T) {
	r, err := NewResolver(&ast.File{})
	require.NoError(t, err)

	cases := []struct {
		Name     string
		List     *ast.IntegerList
		Set      *intervals.Set
		HasError bool
	}{
		{
			Name: "basic",
			List: ast.NewIntegerList(ast.NewIntegerRangeLiteral(4, 5)),
			Set:  intervals.NewSet(intervals.Range(4, 5)),
		},
		{
			Name: "single",
			List: ast.NewIntegerList(ast.NewIntegerRangeSingleLiteral(42)),
			Set:  intervals.NewSet(intervals.Single(42)),
		},
		{
			Name: "multi",
			List: ast.NewIntegerList(
				ast.NewIntegerRangeLiteral(1, 10),
				ast.NewIntegerRangeLiteral(100, 1000),
			),
			Set: intervals.NewSet(
				intervals.Range(1, 10),
				intervals.Range(100, 1000),
			),
		},
		{
			Name: "overlaps",
			List: ast.NewIntegerList(
				ast.NewIntegerRangeLiteral(1, 10),
				ast.NewIntegerRangeLiteral(5, 14),
			),
			HasError: true,
		},
		{
			Name: "lownil",
			List: &ast.IntegerList{
				Ranges: []*ast.IntegerRange{
					&ast.IntegerRange{Low: nil},
				},
			},
			HasError: true,
		},
		{
			Name: "highbadtype",
			List: &ast.IntegerList{
				Ranges: []*ast.IntegerRange{
					&ast.IntegerRange{
						Low:  &ast.IntegerLiteral{Value: 3},
						High: &ast.File{},
					},
				},
			},
			HasError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			s, err := r.Intervals(c.List)
			assert.Equal(t, c.Set, s)
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
