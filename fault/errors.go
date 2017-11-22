// Package fault defines trunnel error types.
package fault

import "fmt"

// UnexpectedType is raised when a function receives an argument of an
// unexpected type. This could happen when visiting a malformed AST, for example.
type UnexpectedType struct {
	t interface{}
}

// NewUnexpectedType builds an UnexpectedType error for the object t.
func NewUnexpectedType(t interface{}) UnexpectedType {
	return UnexpectedType{t: t}
}

func (u UnexpectedType) Error() string {
	return fmt.Sprintf("unexpected type %T", u.t)
}
