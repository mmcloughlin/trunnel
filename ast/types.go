// Package ast defines types used to represent syntax trees for trunnel files.
package ast

// File represents a complete trunnel file.
type File struct {
	Constants []*Constant
	Structs   []*Struct
	Pragmas   []*Pragma
}

// Constant is a constant declaration.
type Constant struct {
	Name  string
	Value int64
}

// Struct is a struct declaration.
type Struct struct {
	Name    string
	Members []StructMember
}

// StructMember is a field in a struct definition.
type StructMember interface{}

// IntegerMember is an integer struct member.
type IntegerMember struct {
	Type       *IntType
	Name       string
	Constraint *IntegerList
}

type ArrayBase interface{}

type FixedArrayMember struct {
	Base ArrayBase
	Name string
	Size Integer
}

// CharType represents the character type.
type CharType struct{}

// IntType represents an integer type (u8, u16, u32 and u64).
type IntType struct {
	Size int
}

// Possible IntTypes.
var (
	U8  = &IntType{Size: 8}
	U16 = &IntType{Size: 16}
	U32 = &IntType{Size: 32}
	U64 = &IntType{Size: 64}
)

// Pragma represents a directive to trunnel.
type Pragma struct {
	Type    string
	Options []string
}

// Integer specifies an integer (either directly or via a constant reference).
type Integer interface{}

// IntegerConstRef specifies an integer via a reference to a constant.
type IntegerConstRef struct {
	Name string
}

// IntegerLiteral specifies an integer directly.
type IntegerLiteral struct {
	Value int64
}

// IntegerRange represents a range of integers.
type IntegerRange struct {
	Low  Integer
	High Integer
}

// IntegerList specifies a set of integers.
type IntegerList struct {
	Ranges []*IntegerRange
}
