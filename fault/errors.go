// Package fault defines trunnel error types.
package fault

import "github.com/pkg/errors"

// ErrNotImplemented indicates a trunnel feature is not implemented.
var ErrNotImplemented = errors.New("not implemented")

// UnexpectedType is raised when a function receives an argument of an
// unexpected type. This could happen when visiting a malformed AST, for example.
type UnexpectedType struct {
	error
}

// NewUnexpectedType builds an UnexpectedType error for the object t.
func NewUnexpectedType(t interface{}) UnexpectedType {
	return UnexpectedType{
		error: errors.Errorf("unexpected type %T", t),
	}
}
