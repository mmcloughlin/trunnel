package inspect

import (
	"testing"

	"github.com/mmcloughlin/trunnel/ast"
	"github.com/mmcloughlin/trunnel/internal/intervals"
	"github.com/mmcloughlin/trunnel/parse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBranchesNoDefault(t *testing.T) {
	b := BuildBranches(t, `struct basic {
		u8 tag;
		union u[tag] {
			1,100: u8 r; u8 g; u8 b;
			2..17, 42: u16 y; u8 m; u8 d;
		};
		u16 right_after_the_union;
	};`)

	require.Len(t, b.branches, 2)
	assert.Equal(t, "1,100", b.branches[0].Set.String())
	assert.Equal(t, "2-17,42", b.branches[1].Set.String())
}

func TestNewBranchesDefault(t *testing.T) {
	b := BuildBranches(t, `struct basic {
		u8 tag;
		union u[tag] {
			0..0xf: u32 a;
            0xf4..0xff: u32 b;
            default: u32 c;
		};
		u16 right_after_the_union;
	};`)

	require.Len(t, b.branches, 3)
	assert.Equal(t, "0-15", b.branches[0].Set.String())
	assert.Equal(t, "244-255", b.branches[1].Set.String())
	assert.Equal(t, "16-243", b.branches[2].Set.String())
}

func TestNewBranchesBadTagRef(t *testing.T) {
	_, err := BuildBranchesWithError(t, `struct basic {
		union u[tagdoesnotexist] {
			0..0xf: u32 a;
            0xf4..0xff: u32 b;
            default: u32 c;
		};
	};`)
	assert.EqualError(t, err, "could not resolve reference")
}

func TestNewBranchesBadIntervals(t *testing.T) {
	_, err := BuildBranchesWithError(t, `struct basic {
        u8 tag;
		union u[tag] {
			0..IDK: u32 a;
		};
	};`)
	assert.EqualError(t, err, "constant undefined")
}

func TestBranchesLookupAndAll(t *testing.T) {
	branches := []Branch{
		{Set: intervals.NewSet(intervals.Range(30, 45), intervals.Range(100, 300))},
		{Set: intervals.NewSet(intervals.Range(10, 20), intervals.Single(1001))},
	}
	b := &Branches{branches: branches}

	assert.Equal(t, branches, b.All())

	cases := []struct {
		X      int64
		Branch Branch
		OK     bool
	}{
		{40, branches[0], true},
		{45, branches[0], true},
		{200, branches[0], true},
		{10, branches[1], true},
		{15, branches[1], true},
		{1001, branches[1], true},
		{10000, Branch{}, false},
		{46, Branch{}, false},
	}
	for _, c := range cases {
		branch, ok := b.Lookup(c.X)
		assert.Equal(t, c.Branch, branch)
		assert.Equal(t, c.OK, ok)
	}
}

func BuildBranches(t *testing.T, code string) *Branches {
	b, err := BuildBranchesWithError(t, code)
	require.NoError(t, err)
	return b
}

func BuildBranchesWithError(t *testing.T, code string) (*Branches, error) {
	f, err := parse.String(code)
	require.NoError(t, err)
	r, err := NewResolver(f)
	require.NoError(t, err)
	s, ok := r.Struct("basic")
	require.True(t, ok)
	u := lookupUnion(s, "u")
	require.NotNil(t, u)
	return NewBranches(r, s, u)
}

func lookupUnion(s *ast.Struct, n string) *ast.UnionMember {
	for _, m := range s.Members {
		if u, ok := m.(*ast.UnionMember); ok {
			return u
		}
	}
	return nil
}
