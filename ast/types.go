// Package ast defines types used to represent syntax trees for trunnel files.
package ast

// File represents a complete trunnel file.
type File struct {
	Constants []*Constant
	Contexts  []*Context
	Structs   []*Struct
	Pragmas   []*Pragma
}

// Declarations
// -----------------------------------------------------------------------------

// Constant is a constant declaration.
type Constant struct {
	Name  string
	Value int64
}

// Context is a context declaration.
type Context struct {
	Name    string
	Members []*Field
}

// Struct is a struct declaration.
type Struct struct {
	Name     string
	Contexts []string
	Members  []Member // nil for extern struct
}

// Extern returns whether the struct declaration is external.
func (s Struct) Extern() bool {
	return s.Members == nil
}

// Pragma represents a directive to trunnel.
type Pragma struct {
	Type    string
	Options []string
}

// Struct Members
// -----------------------------------------------------------------------------

// Member is a field in a struct definition.
type Member interface{}

// Field is a data field in a struct.
type Field struct {
	Name string
	Type Type
}

// UnionMember is a union member of a struct.
type UnionMember struct {
	Name   string
	Tag    *IDRef
	Length LengthConstraint
	Cases  []*UnionCase
}

// EOS signals "end of struct".
type EOS struct{}

// Types
// -----------------------------------------------------------------------------

// Type is a type.
type Type interface{}

// IntType represents an integer type (u8, u16, u32 and u64).
type IntType struct {
	Size       uint
	Constraint *IntegerList
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

// Ptr signals a request to store a pointer to a location within a struct.
type Ptr struct{}

// NulTermString is a NUL-terminated string type.
type NulTermString struct{}

// CharType represents the character type.
type CharType struct{}

// ArrayBase is a type that can be stored in an array.
type ArrayBase interface{}

// FixedArrayMember is a fixed-length array.
type FixedArrayMember struct {
	Base ArrayBase
	Size Integer
}

// VarArrayMember is a variable-length array.
type VarArrayMember struct {
	Base       ArrayBase
	Constraint LengthConstraint // nil means remainder
}

// Unions
// -----------------------------------------------------------------------------

// UnionCase is a case in a union.
type UnionCase struct {
	Case    *IntegerList // nil is the default case
	Members []Member
}

// Fail directive for a union case.
type Fail struct{}

// Ignore directive in a union case.
type Ignore struct{}

// Other
// -----------------------------------------------------------------------------

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
