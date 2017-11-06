package parser

import (
	"strings"
	"testing"

	"github.com/mmcloughlin/trunnel/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ParseString(s string, opts ...Option) (interface{}, error) {
	return ParseReader("", strings.NewReader(s))
}

func TestConstant(t *testing.T) {
	cases := map[string]string{
		"decimal": `const CONST_ID = 123;`,
		"hex":     `const CONST_ID = 0x7b;`,
		"octal":   `const CONST_ID = 0173;`,
	}
	expect := &ast.File{
		Constants: []*ast.Constant{
			{
				Name:  "CONST_ID",
				Value: 123,
			},
		},
	}
	for name, src := range cases {
		t.Run(name, func(t *testing.T) {
			f, err := ParseString(src)
			require.NoError(t, err)
			assert.Equal(t, expect, f)
		})
	}
}

func TestStructBasic(t *testing.T) {
	src := `struct rgb { u8 r; u8 g; u8 b; }`
	expect := &ast.File{
		Structs: []*ast.Struct{
			{
				Name: "rgb",
				Members: []ast.StructMember{
					&ast.IntegerMember{Type: ast.U8, Name: "r"},
					&ast.IntegerMember{Type: ast.U8, Name: "g"},
					&ast.IntegerMember{Type: ast.U8, Name: "b"},
				},
			},
		},
	}
	f, err := ParseString(src)
	require.NoError(t, err)
	assert.Equal(t, expect, f)
}

func TestPragma(t *testing.T) {
	cases := []struct {
		Name    string
		Code    string
		Type    string
		Options []string
	}{
		{
			"single",
			"trunnel options ident1;",
			"options",
			[]string{"ident1"},
		},
		{
			"multi",
			"trunnel special ident1, ident2\t,    ident3   ,ident4   ;",
			"special",
			[]string{"ident1", "ident2", "ident3", "ident4"},
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			f, err := ParseString(c.Code)
			require.NoError(t, err)
			expect := &ast.File{
				Pragmas: []*ast.Pragma{
					{Type: c.Type, Options: c.Options},
				},
			}
			assert.Equal(t, expect, f)
		})
	}
}

func TestExample(t *testing.T) {
	_, err := ParseFile("testdata/example.trunnel")
	require.NoError(t, err)
}
