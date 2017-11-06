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
	cases := []struct {
		Name   string
		Code   string
		Expect int64
	}{
		{"decimal", `const CONST_ID = 123;`, 123},
		{"hex", `const CONST_ID = 0x7b;`, 123},
		{"octal", `const CONST_ID = 0173;`, 123},
		{"zero_decimal", `const CONST_ID = 0;`, 0},
		{"zero_hex", `const CONST_ID = 0x00;`, 0},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			f, err := ParseString(c.Code)
			require.NoError(t, err)
			expect := &ast.File{
				Constants: []*ast.Constant{
					{
						Name:  "CONST_ID",
						Value: c.Expect,
					},
				},
			}
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

func TestIntegerMember(t *testing.T) {
	s := `
	struct int_constraints {
		u8 version_num IN [ 4, 5, 6 ];
		u16 length IN [ 0..16384 ];
		u16 length2 IN [ 0..MAX_LEN ];
		u8 version_num2 IN [ 1, 2, 4..6, 9..128 ];
	};
	`
	expect := &ast.File{
		Structs: []*ast.Struct{
			{
				Name: "int_constraints",
				Members: []ast.StructMember{
					&ast.IntegerMember{
						Type: ast.U8,
						Name: "version_num",
						Constraint: &ast.IntegerList{
							Ranges: []*ast.IntegerRange{
								{Low: &ast.IntegerLiteral{Value: 4}},
								{Low: &ast.IntegerLiteral{Value: 5}},
								{Low: &ast.IntegerLiteral{Value: 6}},
							},
						},
					},
					&ast.IntegerMember{
						Type: ast.U16,
						Name: "length",
						Constraint: &ast.IntegerList{
							Ranges: []*ast.IntegerRange{
								{
									Low:  &ast.IntegerLiteral{Value: 0},
									High: &ast.IntegerLiteral{Value: 16384},
								},
							},
						},
					},
					&ast.IntegerMember{
						Type: ast.U16,
						Name: "length2",
						Constraint: &ast.IntegerList{
							Ranges: []*ast.IntegerRange{
								{
									Low:  &ast.IntegerLiteral{},
									High: &ast.IntegerConstRef{Name: "MAX_LEN"},
								},
							},
						},
					},
					&ast.IntegerMember{
						Type: ast.U8,
						Name: "version_num2",
						Constraint: &ast.IntegerList{
							Ranges: []*ast.IntegerRange{
								{Low: &ast.IntegerLiteral{Value: 1}},
								{Low: &ast.IntegerLiteral{Value: 2}},
								{Low: &ast.IntegerLiteral{Value: 4}, High: &ast.IntegerLiteral{Value: 6}},
								{Low: &ast.IntegerLiteral{Value: 9}, High: &ast.IntegerLiteral{Value: 128}},
							},
						},
					},
				},
			},
		},
	}
	f, err := ParseString(s)
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
