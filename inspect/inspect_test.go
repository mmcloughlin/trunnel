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
