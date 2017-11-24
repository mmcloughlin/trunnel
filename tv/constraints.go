package tv

import (
	"github.com/pkg/errors"

	"github.com/mmcloughlin/trunnel/ast"
)

// Constraints records fixed values for struct/context fields.
type Constraints map[string]map[string]int64

// NewConstraints builds an empty set of constraints.
func NewConstraints() Constraints {
	return Constraints{}
}

// Lookup returns the value of the constraint on s.k.
func (c Constraints) Lookup(s, k string) (int64, bool) {
	values, ok := c[s]
	if !ok {
		return 0, false
	}
	v, ok := values[k]
	return v, ok
}

// LookupLocal is a convenience for looking up in the local scope "".
func (c Constraints) LookupLocal(k string) (int64, bool) {
	return c.Lookup("", k)
}

// LookupRef is a convenience for looking up an AST IDRef.
func (c Constraints) LookupRef(r *ast.IDRef) (int64, bool) {
	return c.Lookup(r.Scope, r.Name)
}

// Set sets the value of s.k.
func (c Constraints) Set(s, k string, v int64) error {
	if _, exists := c[s]; !exists {
		c[s] = map[string]int64{}
	}
	if u, exists := c[s][k]; exists && u != v {
		return errors.New("conflicting constraint")
	}
	c[s][k] = v
	return nil
}

// SetRef is a convenience for setting the value of an AST IDRef.
func (c Constraints) SetRef(r *ast.IDRef, v int64) error {
	return c.Set(r.Scope, r.Name, v)
}

// LookupOrCreate looks up s.k and returns the value if it exists. Otherwise the
// constraint is set to v and returned.
func (c Constraints) LookupOrCreate(s, k string, v int64) int64 {
	if u, ok := c.Lookup(s, k); ok {
		return u
	}
	if err := c.Set(s, k, v); err != nil {
		panic(err) // should not happen, we already checked if it exists
	}
	return v
}

// LookupOrCreateRef is a convenience for LookupOrCreate with an AST IDRef.
func (c Constraints) LookupOrCreateRef(r *ast.IDRef, v int64) int64 {
	return c.LookupOrCreate(r.Scope, r.Name, v)
}

// ClearScope deletes all constraints in scope s.
func (c Constraints) ClearScope(s string) {
	delete(c, s)
}

// ClearLocal deletes all constraints in the local scope.
func (c Constraints) ClearLocal() {
	c.ClearScope("")
}

// Update applies all constraints in d to c.
func (c Constraints) Update(d Constraints) error {
	for s, values := range d {
		for k, v := range values {
			if err := c.Set(s, k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

// Clone returns a deep copy of c.
func (c Constraints) Clone() Constraints {
	clone := NewConstraints()
	if err := clone.Update(c); err != nil {
		panic(err) // theoretically impossible
	}
	return clone
}

// CloneGlobal clones all constraints apart from the local ones. It is a
// convenience for Clone followed by ClearLocal.
func (c Constraints) CloneGlobal() Constraints {
	g := c.Clone()
	g.ClearLocal()
	return g
}

// Merge builds a new set of constraints by merging c and other. Errors on conflict.
func (c Constraints) Merge(other Constraints) (Constraints, error) {
	m := c.Clone()
	if err := m.Update(other); err != nil {
		return nil, err
	}
	return m, nil
}
