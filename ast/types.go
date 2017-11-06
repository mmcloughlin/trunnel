// Package ast defines types used to represent syntax trees for trunnel files.
package ast

type File struct {
	Constants []*Constant
}

type Constant struct {
	Name  string
	Value int64
}

type IntType struct {
	Size int
}
