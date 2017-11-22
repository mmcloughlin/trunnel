// Package fault defines trunnel error types.
package fault

import "fmt"

type UnexpectedType struct {
	t interface{}
}

func NewUnexpectedType(t interface{}) UnexpectedType {
	return UnexpectedType{t: t}
}

func (u UnexpectedType) Error() string {
	return fmt.Sprintf("unexpected type %T", u.t)
}
