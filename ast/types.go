// Package ast defines types used to represent syntax trees for trunnel files.
package ast

type File struct {
	Constants []*Constant
	Structs   []*Struct
}

type Constant struct {
	Name  string
	Value int64
}

type Struct struct {
	Name    string
	Members []StructMember
}

type StructMember interface{}

type IntegerMember struct {
	Type *IntType
	Name string
}

type IntType struct {
	Size int
}

var (
	U8  = &IntType{Size: 8}
	U16 = &IntType{Size: 16}
	U32 = &IntType{Size: 32}
	U64 = &IntType{Size: 64}
)
