package inspect

import (
	"github.com/mmcloughlin/trunnel/ast"
	"github.com/mmcloughlin/trunnel/internal/intervals"
)

// Branch is a case of a union.
type Branch struct {
	Set  *intervals.Set
	Case *ast.UnionCase
}

// Branches represents the branches of a union.
type Branches struct {
	branches []Branch
}

// NewBranches builds a Branches object from a union.
func NewBranches(r *Resolver, s *ast.Struct, u *ast.UnionMember) (*Branches, error) {
	t, err := r.IntType(u.Tag, s)
	if err != nil {
		return nil, err
	}

	dflt := Branch{
		Set: intervals.IntType(t.Size),
	}

	b := &Branches{}
	for _, c := range u.Cases {
		if c.Case == nil {
			dflt.Case = c
			continue
		}

		s, err := r.Intervals(c.Case)
		if err != nil {
			return nil, err
		}

		b.branches = append(b.branches, Branch{
			Set:  s,
			Case: c,
		})

		dflt.Set.Subtract(s)
	}

	if dflt.Case != nil {
		b.branches = append(b.branches, dflt)
	}

	return b, nil
}

// Lookup fetches the branch x falls into.
func (b *Branches) Lookup(x int64) (Branch, bool) {
	for _, branch := range b.branches {
		if branch.Set.Contains(uint64(x)) { // XXX cast
			return branch, true
		}
	}
	return Branch{}, false
}
