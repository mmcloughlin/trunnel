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
}

type IntType struct {
	Size int
}
