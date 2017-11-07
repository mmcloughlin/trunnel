package parser

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmcloughlin/trunnel/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ParseString(s string, opts ...Option) (interface{}, error) {
	return ParseReader("", strings.NewReader(s), opts...)
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
				Members: []ast.Member{
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

func TestIntTypes(t *testing.T) {
	src := `struct s { u8 a; u16 b; u32 c; u64 d; }`
	expect := &ast.File{
		Structs: []*ast.Struct{
			{
				Name: "s",
				Members: []ast.Member{
					&ast.IntegerMember{Type: ast.U8, Name: "a"},
					&ast.IntegerMember{Type: ast.U16, Name: "b"},
					&ast.IntegerMember{Type: ast.U32, Name: "c"},
					&ast.IntegerMember{Type: ast.U64, Name: "d"},
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
				Members: []ast.Member{
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

func TestNulTermString(t *testing.T) {
	src := `struct nul_term_string {
		nulterm str;
	};
	`
	expect := &ast.File{
		Structs: []*ast.Struct{
			{
				Name: "nul_term_string",
				Members: []ast.Member{
					&ast.NulTermString{Name: "str"},
				},
			},
		},
	}
	f, err := ParseString(src)
	require.NoError(t, err)
	assert.Equal(t, expect, f)
}

func TestNestedStructs(t *testing.T) {
	src := `
	struct rgb { u8 r; u8 g; u8 b; };
	struct outer {
		struct rgb color;
		struct inner { u8 a; u64 b; } c;
	};
	`
	expect := &ast.File{
		Structs: []*ast.Struct{
			{
				Name: "rgb",
				Members: []ast.Member{
					&ast.IntegerMember{Type: ast.U8, Name: "r"},
					&ast.IntegerMember{Type: ast.U8, Name: "g"},
					&ast.IntegerMember{Type: ast.U8, Name: "b"},
				},
			},
			{
				Name: "outer",
				Members: []ast.Member{
					&ast.StructMember{Name: "color", Ref: &ast.StructRef{Name: "rgb"}},
					&ast.StructMember{Name: "c", Ref: &ast.StructRef{Name: "inner"}},
				},
			},
			{
				Name: "inner",
				Members: []ast.Member{
					&ast.IntegerMember{Type: ast.U8, Name: "a"},
					&ast.IntegerMember{Type: ast.U64, Name: "b"},
				},
			},
		},
	}
	f, err := ParseString(src)
	require.NoError(t, err)
	assert.Equal(t, expect, f)
}

func TestFixedArraySimple(t *testing.T) {
	src := `struct fixed_arrays {
		u8 a[8];
		u32 b[SIZE];
		char s[13];
	}`
	expect := &ast.File{
		Structs: []*ast.Struct{
			{
				Name: "fixed_arrays",
				Members: []ast.Member{
					&ast.FixedArrayMember{
						Base: ast.U8,
						Name: "a",
						Size: &ast.IntegerLiteral{Value: 8},
					},
					&ast.FixedArrayMember{
						Base: ast.U32,
						Name: "b",
						Size: &ast.IntegerConstRef{Name: "SIZE"},
					},
					&ast.FixedArrayMember{
						Base: &ast.CharType{},
						Name: "s",
						Size: &ast.IntegerLiteral{Value: 13},
					},
				},
			},
		},
	}
	f, err := ParseString(src)
	require.NoError(t, err)
	assert.Equal(t, expect, f)
}

func TestFixedArrayStructs(t *testing.T) {
	src := `struct fixed_array_structs {
		struct another x[3];
		struct inner {
			u8 a;
			u32 b;
		} y[7];
	}`
	expect := &ast.File{
		Structs: []*ast.Struct{
			{
				Name: "fixed_array_structs",
				Members: []ast.Member{
					&ast.FixedArrayMember{
						Base: &ast.StructRef{Name: "another"},
						Name: "x",
						Size: &ast.IntegerLiteral{Value: 3},
					},
					&ast.FixedArrayMember{
						Base: &ast.StructRef{Name: "inner"},
						Name: "y",
						Size: &ast.IntegerLiteral{Value: 7},
					},
				},
			},
			{
				Name: "inner",
				Members: []ast.Member{
					&ast.IntegerMember{Type: ast.U8, Name: "a"},
					&ast.IntegerMember{Type: ast.U32, Name: "b"},
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

func TestValidFiles(t *testing.T) {
	filenames, err := filepath.Glob("testdata/valid/*.trunnel")
	require.NoError(t, err)
	for _, filename := range filenames {
		t.Run(filename, func(t *testing.T) {
			_, err := ParseFile(filename)
			assert.NoError(t, err)
		})
	}
}

func TestFailingFiles(t *testing.T) {
	filenames, err := filepath.Glob("testdata/failing/*.trunnel")
	require.NoError(t, err)
	for _, filename := range filenames {
		t.Run(filename, func(t *testing.T) {
			_, err := ParseFile(filename)
			assert.Error(t, err)
		})
	}
}

// TestOptions is primarily provided for test coverage of the generated Option
// functions.
func TestOptions(t *testing.T) {
	opts := []Option{
		AllowInvalidUTF8(false),
		Debug(false),
		Entrypoint(""),
		MaxExpressions(0),
		Memoize(false),
		Recover(true),
		GlobalStore("foo", "baz"),
		InitState("blah", 42),
		Statistics(&Stats{}, "hmm"),
	}
	_, err := ParseString("const A = 1337;", opts...)
	assert.NoError(t, err)
}
