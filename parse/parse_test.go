package parse

import (
	"testing"

	"github.com/mmcloughlin/trunnel/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmpty(t *testing.T) {
	empty := map[string]string{
		"blank":  "",
		"spaces": "    ",
		"tab":    "\t",
		"single_line_comment": "// nothing here",
		"multi_line_comment":  "   /* or here*/\t",
	}
	for n, src := range empty {
		t.Run(n, func(t *testing.T) {
			f, err := String(src)
			require.NoError(t, err)
			assert.Equal(t, &ast.File{}, f)
		})
	}
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
			f, err := String(c.Code)
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
					&ast.Field{Type: ast.U8, Name: "r"},
					&ast.Field{Type: ast.U8, Name: "g"},
					&ast.Field{Type: ast.U8, Name: "b"},
				},
			},
		},
	}
	f, err := String(src)
	require.NoError(t, err)
	assert.Equal(t, expect, f)
}

func TestExtern(t *testing.T) {
	src := "extern struct rgb;"
	expect := &ast.File{
		Extern: []*ast.ExternStruct{
			&ast.ExternStruct{Name: "rgb"},
		},
	}
	f, err := String(src)
	require.NoError(t, err)
	assert.Equal(t, expect, f)
}

func TestExternContexts(t *testing.T) {
	src := "extern struct rgb with context a,b,c;"
	expect := &ast.File{
		Extern: []*ast.ExternStruct{
			&ast.ExternStruct{
				Name:     "rgb",
				Contexts: []string{"a", "b", "c"},
			},
		},
	}
	f, err := String(src)
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
					&ast.Field{Type: ast.U8, Name: "a"},
					&ast.Field{Type: ast.U16, Name: "b"},
					&ast.Field{Type: ast.U32, Name: "c"},
					&ast.Field{Type: ast.U64, Name: "d"},
				},
			},
		},
	}
	f, err := String(src)
	require.NoError(t, err)
	assert.Equal(t, expect, f)
}

func TestField(t *testing.T) {
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
					&ast.Field{
						Name: "version_num",
						Type: &ast.IntType{
							Size: 8,
							Constraint: &ast.IntegerList{
								Ranges: []*ast.IntegerRange{
									{Low: &ast.IntegerLiteral{Value: 4}},
									{Low: &ast.IntegerLiteral{Value: 5}},
									{Low: &ast.IntegerLiteral{Value: 6}},
								},
							},
						},
					},
					&ast.Field{
						Name: "length",
						Type: &ast.IntType{
							Size: 16,
							Constraint: &ast.IntegerList{
								Ranges: []*ast.IntegerRange{
									{
										Low:  &ast.IntegerLiteral{Value: 0},
										High: &ast.IntegerLiteral{Value: 16384},
									},
								},
							},
						},
					},
					&ast.Field{
						Name: "length2",
						Type: &ast.IntType{
							Size: 16,
							Constraint: &ast.IntegerList{
								Ranges: []*ast.IntegerRange{
									{
										Low:  &ast.IntegerLiteral{},
										High: &ast.IntegerConstRef{Name: "MAX_LEN"},
									},
								},
							},
						},
					},
					&ast.Field{
						Name: "version_num2",
						Type: &ast.IntType{
							Size: 8,
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
		},
	}
	f, err := String(s)
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
					&ast.Field{
						Name: "str",
						Type: &ast.NulTermString{},
					},
				},
			},
		},
	}
	f, err := String(src)
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
					&ast.Field{Type: ast.U8, Name: "r"},
					&ast.Field{Type: ast.U8, Name: "g"},
					&ast.Field{Type: ast.U8, Name: "b"},
				},
			},
			{
				Name: "outer",
				Members: []ast.Member{
					&ast.Field{Name: "color", Type: &ast.StructRef{Name: "rgb"}},
					&ast.Field{Name: "c", Type: &ast.StructRef{Name: "inner"}},
				},
			},
			{
				Name: "inner",
				Members: []ast.Member{
					&ast.Field{Type: ast.U8, Name: "a"},
					&ast.Field{Type: ast.U64, Name: "b"},
				},
			},
		},
	}
	f, err := String(src)
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
					&ast.Field{
						Name: "a",
						Type: &ast.FixedArrayMember{
							Base: ast.U8,
							Size: &ast.IntegerLiteral{Value: 8},
						},
					},
					&ast.Field{
						Name: "b",
						Type: &ast.FixedArrayMember{
							Base: ast.U32,
							Size: &ast.IntegerConstRef{Name: "SIZE"},
						},
					},
					&ast.Field{
						Name: "s",
						Type: &ast.FixedArrayMember{
							Base: &ast.CharType{},
							Size: &ast.IntegerLiteral{Value: 13},
						},
					},
				},
			},
		},
	}
	f, err := String(src)
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
					&ast.Field{
						Name: "x",
						Type: &ast.FixedArrayMember{
							Base: &ast.StructRef{Name: "another"},
							Size: &ast.IntegerLiteral{Value: 3},
						},
					},
					&ast.Field{
						Name: "y",
						Type: &ast.FixedArrayMember{
							Base: &ast.StructRef{Name: "inner"},
							Size: &ast.IntegerLiteral{Value: 7},
						},
					},
				},
			},
			{
				Name: "inner",
				Members: []ast.Member{
					&ast.Field{Type: ast.U8, Name: "a"},
					&ast.Field{Type: ast.U32, Name: "b"},
				},
			},
		},
	}
	f, err := String(src)
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
					&ast.Field{Type: ast.U16, Name: "length"},
					&ast.Field{
						Name: "bytes",
						Type: &ast.VarArrayMember{
							Base:       ast.U8,
							Constraint: &ast.IDRef{Name: "length"},
						},
					},
				},
			},
		},
	}
	f, err := String(src)
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
					&ast.Field{Type: ast.U8, Name: "x"},
					&ast.Field{
						Name: "rest",
						Type: &ast.VarArrayMember{
							Base:       ast.U8,
							Constraint: nil,
						},
					},
				},
			},
		},
	}
	f, err := String(src)
	require.NoError(t, err)
	assert.Equal(t, expect, f)
}

func TestEOS(t *testing.T) {
	src := ` struct fourbytes {
		u16 x;
		u16 y;
		eos;
	}`
	expect := &ast.File{
		Structs: []*ast.Struct{
			{
				Name: "fourbytes",
				Members: []ast.Member{
					&ast.Field{Type: ast.U16, Name: "x"},
					&ast.Field{Type: ast.U16, Name: "y"},
					&ast.EOS{},
				},
			},
		},
	}
	f, err := String(src)
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
					&ast.Field{Type: ast.U8, Name: "hostname_len"},
					&ast.Field{
						Name: "hostname",
						Type: &ast.VarArrayMember{
							Base:       &ast.CharType{},
							Constraint: &ast.IDRef{Name: "hostname_len"},
						},
					},
				},
			},
		},
	}
	f, err := String(src)
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
					&ast.Field{
						Name: "salt",
						Type: &ast.FixedArrayMember{
							Base: ast.U8,
							Size: &ast.IntegerLiteral{Value: 16},
						},
					},
					&ast.Field{
						Name: "message",
						Type: &ast.VarArrayMember{
							Base: ast.U8,
							Constraint: &ast.Leftover{
								Num: &ast.IntegerLiteral{Value: 32},
							},
						},
					},
					&ast.Field{
						Name: "mac",
						Type: &ast.FixedArrayMember{
							Base: ast.U8,
							Size: &ast.IntegerLiteral{Value: 32},
						},
					},
				},
			},
		},
	}
	f, err := String(src)
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
					&ast.Field{Type: ast.U8, Name: "tag"},
					&ast.UnionMember{
						Name: "addr",
						Tag:  &ast.IDRef{Name: "tag"},
						Cases: []*ast.UnionCase{
							&ast.UnionCase{
								Case: ast.NewIntegerList(ast.NewIntegerRangeSingleLiteral(4)),
								Members: []ast.Member{
									&ast.Field{Type: ast.U32, Name: "ipv4_addr"},
								},
							},
							&ast.UnionCase{
								Case: ast.NewIntegerList(ast.NewIntegerRangeSingleLiteral(5)),
							},
							&ast.UnionCase{
								Case: ast.NewIntegerList(ast.NewIntegerRangeSingleLiteral(6)),
								Members: []ast.Member{
									&ast.Field{
										Name: "ipv6_addr",
										Type: &ast.FixedArrayMember{
											Base: ast.U8,
											Size: &ast.IntegerLiteral{Value: 16},
										},
									},
								},
							},
							&ast.UnionCase{
								Case: ast.NewIntegerList(
									ast.NewIntegerRangeSingleLiteral(0xf0),
									ast.NewIntegerRangeSingleLiteral(0xf1),
								),
								Members: []ast.Member{
									&ast.Field{Type: ast.U8, Name: "hostname_len"},
									&ast.Field{
										Name: "hostname",
										Type: &ast.VarArrayMember{
											Base:       &ast.CharType{},
											Constraint: &ast.IDRef{Name: "hostname_len"},
										},
									},
								},
							},
							&ast.UnionCase{
								Case: ast.NewIntegerList(
									ast.NewIntegerRangeLiteral(0xf2, 0xff),
								),
								Members: []ast.Member{
									&ast.Field{Name: "ext", Type: &ast.StructRef{Name: "extension"}},
								},
							},
							&ast.UnionCase{
								Case:    nil, // default
								Members: []ast.Member{&ast.Fail{}},
							},
						},
					},
				},
			},
		},
	}
	f, err := String(src)
	require.NoError(t, err)
	assert.Equal(t, expect, f)
}

func TestUnionExtentSpec(t *testing.T) {
	src := `struct union_extent {
		u8 tag;
		u16 length;
		union addr[tag] with length length {
		   7 : ignore;
		   0xEE : u32 ipv4_addr;
		          ...;
		   0xEF : u32 ipv4_addr;
		          u8 remainder[];
		   default: u8 unrecognized[];
		};
	};`
	expect := &ast.File{
		Structs: []*ast.Struct{
			{
				Name: "union_extent",
				Members: []ast.Member{
					&ast.Field{Type: ast.U8, Name: "tag"},
					&ast.Field{Type: ast.U16, Name: "length"},
					&ast.UnionMember{
						Name:   "addr",
						Tag:    &ast.IDRef{Name: "tag"},
						Length: &ast.IDRef{Name: "length"},
						Cases: []*ast.UnionCase{
							&ast.UnionCase{
								Case:    ast.NewIntegerList(ast.NewIntegerRangeSingleLiteral(7)),
								Members: []ast.Member{&ast.Ignore{}},
							},
							&ast.UnionCase{
								Case: ast.NewIntegerList(ast.NewIntegerRangeSingleLiteral(0xee)),
								Members: []ast.Member{
									&ast.Field{Type: ast.U32, Name: "ipv4_addr"},
									&ast.Ignore{},
								},
							},
							&ast.UnionCase{
								Case: ast.NewIntegerList(ast.NewIntegerRangeSingleLiteral(0xef)),
								Members: []ast.Member{
									&ast.Field{Type: ast.U32, Name: "ipv4_addr"},
									&ast.Field{
										Name: "remainder",
										Type: &ast.VarArrayMember{
											Base:       ast.U8,
											Constraint: nil,
										},
									},
								},
							},
							&ast.UnionCase{
								Case: nil,
								Members: []ast.Member{
									&ast.Field{
										Name: "unrecognized",
										Type: &ast.VarArrayMember{
											Base:       ast.U8,
											Constraint: nil,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	f, err := String(src)
	require.NoError(t, err)
	assert.Equal(t, expect, f)
}

func TestUnionMembersAfter(t *testing.T) {
	src := `struct encrypted {
	   u8 type;
	   union u[type] with length ..-32 {
	      1: u8 bytes[];
	  2: u8 salt[16];
	     u8 other_bytes[];
	   };
	   u64 data[4];
	}`
	expect := &ast.File{
		Structs: []*ast.Struct{
			{
				Name: "encrypted",
				Members: []ast.Member{
					&ast.Field{Type: ast.U8, Name: "type"},
					&ast.UnionMember{
						Name:   "u",
						Tag:    &ast.IDRef{Name: "type"},
						Length: &ast.Leftover{Num: &ast.IntegerLiteral{Value: 32}},
						Cases: []*ast.UnionCase{
							&ast.UnionCase{
								Case: ast.NewIntegerList(ast.NewIntegerRangeSingleLiteral(1)),
								Members: []ast.Member{
									&ast.Field{
										Name: "bytes",
										Type: &ast.VarArrayMember{
											Base:       ast.U8,
											Constraint: nil,
										},
									},
								},
							},
							&ast.UnionCase{
								Case: ast.NewIntegerList(ast.NewIntegerRangeSingleLiteral(2)),
								Members: []ast.Member{
									&ast.Field{
										Name: "salt",
										Type: &ast.FixedArrayMember{
											Base: ast.U8,
											Size: &ast.IntegerLiteral{Value: 16},
										},
									},
									&ast.Field{
										Name: "other_bytes",
										Type: &ast.VarArrayMember{
											Base:       ast.U8,
											Constraint: nil,
										},
									},
								},
							},
						},
					},
					&ast.Field{
						Name: "data",
						Type: &ast.FixedArrayMember{
							Base: ast.U64,
							Size: &ast.IntegerLiteral{Value: 4},
						},
					},
				},
			},
		},
	}
	f, err := String(src)
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
			f, err := String(c.Code)
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

func TestComments(t *testing.T) {
	src := `struct /* comments can
	be anywhere */ rgb {
		u8 r; /* this is a multi
	line comment that /*should exclude this:
		u8 a;
	(hopefully) */
		u8 // end of line comment
			g; //}
		u8 b;
	}`
	expect := &ast.File{
		Structs: []*ast.Struct{
			{
				Name: "rgb",
				Members: []ast.Member{
					&ast.Field{Type: ast.U8, Name: "r"},
					&ast.Field{Type: ast.U8, Name: "g"},
					&ast.Field{Type: ast.U8, Name: "b"},
				},
			},
		},
	}
	f, err := String(src)
	require.NoError(t, err)
	assert.Equal(t, expect, f)
}

func TestContext(t *testing.T) {
	src := `context ctx { u8 a; u16 b; u32 c; u64 d; }`
	expect := &ast.File{
		Contexts: []*ast.Context{
			{
				Name: "ctx",
				Members: []*ast.Field{
					&ast.Field{Type: ast.U8, Name: "a"},
					&ast.Field{Type: ast.U16, Name: "b"},
					&ast.Field{Type: ast.U32, Name: "c"},
					&ast.Field{Type: ast.U64, Name: "d"},
				},
			},
		},
	}
	f, err := String(src)
	require.NoError(t, err)
	assert.Equal(t, expect, f)
}

// TestContextStructMemberErrors confirms we get errors for member types that
// are valid in structs but not in contexts.
func TestContextStructMemberErrors(t *testing.T) {
	members := map[string]string{
		"fixed_array": "u8 fixed_array[2];",
		"var_array":   "u16 len; u8 var_array[len];",
		"eos":         "u8 a; eos;",
		"remaining":   "u32 a; u8 rest[];",
		"union": `u8 tag; union u[tag] {
			1 : ignore;
			default: fail;
		};`,
		"int_constraint": "u8 x IN [ 42 ];",
	}
	for n, m := range members {
		t.Run(n, func(t *testing.T) {
			_, err := String("struct verify {" + m + "}")
			require.NoError(t, err)
			_, err = String("context ctx {" + m + "}")
			assert.Error(t, err)
		})
	}
}

func TestStructWithContext(t *testing.T) {
	src := `struct encrypted_record with context stream_settings {
		u8 iv[stream_settings.iv_len];
	};`
	expect := &ast.File{
		Structs: []*ast.Struct{
			{
				Name:     "encrypted_record",
				Contexts: []string{"stream_settings"},
				Members: []ast.Member{
					&ast.Field{
						Name: "iv",
						Type: &ast.VarArrayMember{
							Base: ast.U8,
							Constraint: &ast.IDRef{
								Scope: "stream_settings",
								Name:  "iv_len",
							},
						},
					},
				},
			},
		},
	}
	f, err := String(src)
	require.NoError(t, err)
	assert.Equal(t, expect, f)
}

func TestStructWithMultipleContexts(t *testing.T) {
	src := `struct multi with context ctx0,ctx1 , ctx2,      ctx3 {
		u8 x;
	};`
	expect := &ast.File{
		Structs: []*ast.Struct{
			{
				Name:     "multi",
				Contexts: []string{"ctx0", "ctx1", "ctx2", "ctx3"},
				Members: []ast.Member{
					&ast.Field{Type: ast.U8, Name: "x"},
				},
			},
		},
	}
	f, err := String(src)
	require.NoError(t, err)
	assert.Equal(t, expect, f)
}

func TestStructPtr(t *testing.T) {
	src := `struct haspos {
	  nulterm s1;
	  /** Position right after the first NUL. */
	  @ptr pos1;
	  nulterm s2;
	  @ptr pos2;
	  u32 x;
	}`
	expect := &ast.File{
		Structs: []*ast.Struct{
			{
				Name: "haspos",
				Members: []ast.Member{
					&ast.Field{Name: "s1", Type: &ast.NulTermString{}},
					&ast.Field{Name: "pos1", Type: &ast.Ptr{}},
					&ast.Field{Name: "s2", Type: &ast.NulTermString{}},
					&ast.Field{Name: "pos2", Type: &ast.Ptr{}},
					&ast.Field{Type: ast.U32, Name: "x"},
				},
			},
		},
	}
	f, err := String(src)
	require.NoError(t, err)
	assert.Equal(t, expect, f)
}
