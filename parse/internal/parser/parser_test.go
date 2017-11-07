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

func TestVarLengthArray(t *testing.T) {
	src := `struct var_length_array {
		u16 length;
		u8 bytes[length];
	}`
	expect := &ast.File{
		Structs: []*ast.Struct{
			{
				Name: "var_length_array",
				Members: []ast.Member{
					&ast.IntegerMember{Type: ast.U16, Name: "length"},
					&ast.VarArrayMember{
						Base:       ast.U8,
						Name:       "bytes",
						Constraint: &ast.IDRef{Name: "length"},
					},
				},
			},
		},
	}
	f, err := ParseString(src)
	require.NoError(t, err)
	assert.Equal(t, expect, f)
}

func TestRemainderArray(t *testing.T) {
	src := `struct remainder_array {
		u8 x;
		u8 rest[];
	}`
	expect := &ast.File{
		Structs: []*ast.Struct{
			{
				Name: "remainder_array",
				Members: []ast.Member{
					&ast.IntegerMember{Type: ast.U8, Name: "x"},
					&ast.VarArrayMember{
						Base:       ast.U8,
						Name:       "rest",
						Constraint: nil,
					},
				},
			},
		},
	}
	f, err := ParseString(src)
	require.NoError(t, err)
	assert.Equal(t, expect, f)
}

func TestVarLengthString(t *testing.T) {
	src := `struct pascal_string {
		u8 hostname_len;
		char hostname[hostname_len];
	}`
	expect := &ast.File{
		Structs: []*ast.Struct{
			{
				Name: "pascal_string",
				Members: []ast.Member{
					&ast.IntegerMember{Type: ast.U8, Name: "hostname_len"},
					&ast.VarArrayMember{
						Base:       &ast.CharType{},
						Name:       "hostname",
						Constraint: &ast.IDRef{Name: "hostname_len"},
					},
				},
			},
		},
	}
	f, err := ParseString(src)
	require.NoError(t, err)
	assert.Equal(t, expect, f)
}

func TestLeftoverLengthArray(t *testing.T) {
	src := `struct encrypted {
		u8 salt[16];
		u8 message[..-32];
		u8 mac[32];
	}`
	expect := &ast.File{
		Structs: []*ast.Struct{
			{
				Name: "encrypted",
				Members: []ast.Member{
					&ast.FixedArrayMember{
						Base: ast.U8,
						Name: "salt",
						Size: &ast.IntegerLiteral{Value: 16},
					},
					&ast.VarArrayMember{
						Base: ast.U8,
						Name: "message",
						Constraint: &ast.Leftover{
							Num: &ast.IntegerLiteral{Value: 32},
						},
					},
					&ast.FixedArrayMember{
						Base: ast.U8,
						Name: "mac",
						Size: &ast.IntegerLiteral{Value: 32},
					},
				},
			},
		},
	}
	f, err := ParseString(src)
	require.NoError(t, err)
	assert.Equal(t, expect, f)
}

func TestUnion(t *testing.T) {
	src := `struct has_union {
		u8 tag;
		union addr[tag] {
			4 : u32 ipv4_addr;
			5 : ;
			6 : u8 ipv6_addr[16];
			0xf0,0xf1 : u8 hostname_len;
					char hostname[hostname_len];
			0xF2 .. 0xFF : struct extension ext;
			default : fail;
		};
	}`
	expect := &ast.File{
		Structs: []*ast.Struct{
			{
				Name: "has_union",
				Members: []ast.Member{
					&ast.IntegerMember{Type: ast.U8, Name: "tag"},
					&ast.UnionMember{
						Name: "addr",
						Tag:  &ast.IDRef{Name: "tag"},
						Cases: []interface{}{
							&ast.UnionCase{
								Case: ast.NewIntegerList(ast.NewIntegerRangeSingleLiteral(4)),
								Fields: []ast.Member{
									&ast.IntegerMember{Type: ast.U32, Name: "ipv4_addr"},
								},
							},
							&ast.UnionCase{
								Case: ast.NewIntegerList(ast.NewIntegerRangeSingleLiteral(5)),
							},
							&ast.UnionCase{
								Case: ast.NewIntegerList(ast.NewIntegerRangeSingleLiteral(6)),
								Fields: []ast.Member{
									&ast.FixedArrayMember{
										Base: ast.U8,
										Name: "ipv6_addr",
										Size: &ast.IntegerLiteral{Value: 16},
									},
								},
							},
							&ast.UnionCase{
								Case: ast.NewIntegerList(
									ast.NewIntegerRangeSingleLiteral(0xf0),
									ast.NewIntegerRangeSingleLiteral(0xf1),
								),
								Fields: []ast.Member{
									&ast.IntegerMember{Type: ast.U8, Name: "hostname_len"},
									&ast.VarArrayMember{
										Base:       &ast.CharType{},
										Name:       "hostname",
										Constraint: &ast.IDRef{Name: "hostname_len"},
									},
								},
							},
							&ast.UnionCase{
								Case: ast.NewIntegerList(
									ast.NewIntegerRangeLiteral(0xf2, 0xff),
								),
								Fields: []ast.Member{
									&ast.StructMember{Name: "ext", Ref: &ast.StructRef{Name: "extension"}},
								},
							},
							&ast.UnionCase{
								Case:   nil, // default
								Fields: &ast.Fail{},
							},
						},
					},
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
