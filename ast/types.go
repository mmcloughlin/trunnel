// Package ast defines types used to represent syntax trees for trunnel files.
package ast

// File represents a complete trunnel file.
type File struct {
	Constants []*Constant
	Contexts  []*Context
	Structs   []*Struct
	Extern    []*ExternStruct
	Pragmas   []*Pragma
}

// Constant is a constant declaration.
type Constant struct {
	Name  string
	Value int64
}

// Context is a context declaration.
type Context struct {
	Name    string
	Members []*IntegerMember
}

// Struct is a struct declaration.
type Struct struct {
	Name     string
	Contexts []string
	Members  []Member
}

// Member is a field in a struct definition.
type Member interface{}

// IntegerMember is an integer struct member.
type IntegerMember struct {
	Type       *IntType
	Name       string
	Constraint *IntegerList
}

// ArrayBase is a type that can be stored in an array.
type ArrayBase interface{}

// FixedArrayMember is a fixed-length array struct member.
type FixedArrayMember struct {
	Base ArrayBase
	Name string
	Size Integer
}

// VarArrayMember is a variable-length array struct member.
type VarArrayMember struct {
	Base       ArrayBase
	Name       string
	Constraint LengthConstraint // nil means remainder
}

// Ptr signals a request to store a pointer to a location within a struct.
type Ptr struct {
	Name string
}

// EOS signals "end of struct".
type EOS struct{}

// NulTermString is a NUL-terminated string struct member.
type NulTermString struct {
	Name string
}

// StructMember is a struct type member of a struct.
type StructMember struct {
	Ref  *StructRef
	Name string
}

// UnionMember is a union member of a struct.
type UnionMember struct {
	Name   string
	Tag    *IDRef
	Length LengthConstraint
	Cases  []*UnionCase
}

// UnionCase is a case in a union.
type UnionCase struct {
	Case   *IntegerList // nil is the default case
	Fields interface{}
}

// Fail directive for a union case.
type Fail struct{}

// Ignore directive in a union case.
type Ignore struct{}

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

// StructRef represents a reference to a struct type.
type StructRef struct {
	Name string
}

// ExternStruct is a declaration that a Trunnel structure is available
// elsewhere.
type ExternStruct struct {
	Name     string
	Contexts []string
}

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

// NewIntegerRange constructs an IntegerRange from lo to hi.
func NewIntegerRange(lo, hi Integer) *IntegerRange {
	return &IntegerRange{
		Low:  lo,
		High: hi,
	}
}

// NewIntegerRangeLiteral constructs an IntegerRange with literal bounds.
func NewIntegerRangeLiteral(lo, hi int64) *IntegerRange {
	return NewIntegerRange(
		&IntegerLiteral{Value: lo},
		&IntegerLiteral{Value: hi},
	)
}

// NewIntegerRangeSingle constructs an IntegerRange containing just one integer.
func NewIntegerRangeSingle(i Integer) *IntegerRange {
	return NewIntegerRange(i, nil)
}

// NewIntegerRangeSingleLiteral constructs an IntegerRange containing a single
// integer specified with a literal.
func NewIntegerRangeSingleLiteral(v int64) *IntegerRange {
	return NewIntegerRangeSingle(&IntegerLiteral{Value: v})
}

// IntegerList specifies a set of integers.
type IntegerList struct {
	Ranges []*IntegerRange
}

// NewIntegerList constructs an integer list from the given ranges.
func NewIntegerList(ranges ...*IntegerRange) *IntegerList {
	return &IntegerList{
		Ranges: ranges,
	}
}

// LengthConstraint specifies a constraint on the length of a struct member.
type LengthConstraint interface{}

// Leftover is a LengthConstraint which specifies the member occupies all but
// the last Num bytes.
type Leftover struct {
	Num Integer
}

// IDRef is a reference to an identifier, possibly within a scope.
type IDRef struct {
	Scope string
	Name  string
}
