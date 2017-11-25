package tv

import (
	"github.com/mmcloughlin/random"
	"github.com/mmcloughlin/trunnel/ast"
	"github.com/mmcloughlin/trunnel/fault"
	"github.com/mmcloughlin/trunnel/inspect"
	"github.com/pkg/errors"
)

// Vector is a test vector.
type Vector struct {
	Data        []byte
	Constraints Constraints
}

// NewVector builds a test vector with empty constraints.
func NewVector(b []byte) Vector {
	return Vector{
		Data:        b,
		Constraints: NewConstraints(),
	}
}

func cross(a, b []Vector) ([]Vector, error) {
	p := []Vector{}
	for _, u := range a {
		for _, w := range b {
			m, err := u.Constraints.Merge(w.Constraints)
			if err != nil {
				return nil, err
			}
			v := Vector{
				Data:        append(u.Data, w.Data...),
				Constraints: m,
			}
			p = append(p, v)
		}
	}
	return p, nil
}

type generator struct {
	resolver    *inspect.Resolver
	constraints Constraints
	strct       *ast.Struct

	rnd random.Interface
}

// Generate generates a set of test vectors for the types defined in f.
func Generate(f *ast.File, opts ...Option) (map[string][]Vector, error) {
	g := &generator{
		rnd: random.New(),
	}
	for _, opt := range opts {
		opt(g)
	}
	return g.file(f)
}

// Option is an option to control test vector generation.
type Option func(*generator)

// WithRandom sets the random source for test vector generation.
func WithRandom(r random.Interface) Option {
	return func(g *generator) {
		g.rnd = r
	}
}

func (g *generator) init(f *ast.File) (err error) {
	g.resolver, err = inspect.NewResolver(f)
	return
}

func (g *generator) file(f *ast.File) (map[string][]Vector, error) {
	err := g.init(f)
	if err != nil {
		return nil, err
	}

	v := map[string][]Vector{}
	for _, s := range f.Structs {
		if s.Extern() {
			continue
		}
		g.constraints = NewConstraints()
		v[s.Name], err = g.structure(s)
		if err != nil {
			return nil, err
		}
	}

	return v, nil
}

func (g *generator) structure(s *ast.Struct) ([]Vector, error) {
	restore := g.strct
	g.strct = s
	vs := []Vector{
		{
			Data:        []byte{},
			Constraints: g.constraints.CloneGlobal(),
		},
	}
	vs, err := g.members(vs, s.Members)
	if err != nil {
		return nil, err
	}
	g.constraints.ClearLocal()
	for _, v := range vs {
		v.Constraints.ClearLocal()
	}
	g.strct = restore
	return vs, nil
}

func (g *generator) members(vs []Vector, ms []ast.Member) ([]Vector, error) {
	n := len(ms)
	for i := n - 1; i >= 0; i-- {
		extended := []Vector{}
		for _, v := range vs {
			g.constraints = v.Constraints
			mvs, err := g.member(ms[i])
			if err != nil {
				return nil, err
			}

			for _, mv := range mvs {
				mv.Data = append(mv.Data, v.Data...)
				extended = append(extended, mv)
			}
		}
		vs = extended
	}
	return vs, nil
}

func (g *generator) member(m ast.Member) ([]Vector, error) {
	switch m := m.(type) {
	case *ast.Field:
		return g.field(m)
	case *ast.UnionMember:
		return g.union(m)
	case *ast.Ignore:
		return []Vector{
			g.empty(),
			g.vector(g.randbytes(1, 7)),
		}, nil
	case *ast.Fail:
		return []Vector{}, nil
	case *ast.EOS:
		return []Vector{g.empty()}, nil
	default:
		return nil, fault.NewUnexpectedType(m)
	}
}

func (g *generator) field(f *ast.Field) ([]Vector, error) {
	switch t := f.Type.(type) {
	case *ast.IntType:
		return g.intType(f.Name, t)

	case *ast.CharType:
		return g.intType(f.Name, ast.U8)

	case *ast.Ptr:
		return []Vector{g.empty()}, nil

	case *ast.NulTermString:
		return []Vector{g.vector(g.randnulterm(2, 20))}, nil

	case *ast.StructRef:
		s, ok := g.resolver.Struct(t.Name)
		if !ok {
			return nil, errors.New("could not resolve struct name")
		}
		return g.structure(s)

	case *ast.FixedArrayMember:
		return g.array(t.Base, t.Size)

	case *ast.VarArrayMember:
		return g.array(t.Base, t.Constraint)

	default:
		return nil, fault.NewUnexpectedType(t)
	}
}

func (g *generator) intType(name string, t *ast.IntType) ([]Vector, error) {
	var b []byte

	switch x, ok := g.constraints.LookupLocal(name); {
	case ok:
		b = intbytes(x, t.Size)

	case t.Constraint != nil:
		s, err := g.resolver.Intervals(t.Constraint)
		if err != nil {
			return nil, err
		}
		r := s.Random()
		b = intbytes(int64(r), t.Size) // XXX cast

	default:
		b = g.randint(t.Size)
	}

	return []Vector{g.vector(b)}, nil
}

func (g *generator) array(base ast.Type, s ast.LengthConstraint) ([]Vector, error) {
	iv := g.empty()
	v := []Vector{iv}

	var n int64
	switch s := s.(type) {
	case *ast.IntegerConstRef, *ast.IntegerLiteral:
		i, err := g.resolver.Integer(s)
		if err != nil {
			return nil, err
		}
		n = i

	case *ast.IDRef:
		n = iv.Constraints.LookupOrCreateRef(s, int64(g.randbtw(1, 20)))

	case nil:
		n = int64(g.randbtw(1, 20))

	case *ast.Leftover:
		return nil, fault.ErrNotImplemented

	default:
		return nil, fault.NewUnexpectedType(s)
	}

	for i := int64(0); i < n; i++ {
		w, err := g.field(&ast.Field{Type: base})
		if err != nil {
			return nil, err
		}
		v, err = cross(w, v)
		if err != nil {
			return nil, err
		}
	}
	return v, nil
}

func (g *generator) union(u *ast.UnionMember) ([]Vector, error) {
	branches, err := inspect.NewBranches(g.resolver, g.strct, u)
	if err != nil {
		return nil, err
	}

	// has the tag already been set?
	options := branches.All()
	t, ok := g.constraints.LookupRef(u.Tag)
	if ok {
		branch, ok := branches.Lookup(t)
		if !ok {
			return []Vector{}, nil
		}
		options = []inspect.Branch{branch}
	}

	base := g.constraints.Clone()
	results := []Vector{}

	for _, b := range options {
		g.constraints = base.Clone()
		g.constraints.LookupOrCreateRef(u.Tag, int64(b.Set.Random())) // XXX cast
		vs, err := g.members([]Vector{g.empty()}, b.Case.Members)
		if err != nil {
			return nil, err
		}
		results = append(results, vs...)
	}

	// set length constraint
	if u.Length != nil {
		return g.lenconstrain(u.Length, results)
	}

	return results, nil
}

func (g *generator) lenconstrain(c ast.LengthConstraint, vs []Vector) ([]Vector, error) {
	r, ok := c.(*ast.IDRef)
	if !ok {
		return nil, fault.ErrNotImplemented
	}

	results := []Vector{}
	for _, v := range vs {
		n := int64(len(v.Data))
		cst := v.Constraints.Clone()
		m := cst.LookupOrCreateRef(r, n)
		if m != n {
			continue
		}
		results = append(results, Vector{
			Data:        v.Data,
			Constraints: cst,
		})
	}

	return results, nil
}

// empty returns an empty Vector with current constraints.
func (g *generator) empty() Vector {
	return g.vector([]byte{})
}

// vector builds vector with the current constraints.
func (g *generator) vector(b []byte) Vector {
	return Vector{
		Data:        b,
		Constraints: g.constraints,
	}
}

func (g *generator) randint(bits uint) []byte {
	b := make([]byte, bits/8)
	g.randread(b)
	return b
}

// randbtw generates an integer between a and b, inclusive.
func (g *generator) randbtw(a, b int) int {
	return a + g.rnd.Intn(b-a+1)
}

// randbytes returns a random byre array of length between a and b.
func (g *generator) randbytes(a, b int) []byte {
	d := make([]byte, g.randbtw(a, b))
	g.randread(d)
	return d
}

// randread reads random bytes into b from the configured random source.
func (g *generator) randread(b []byte) {
	if _, err := g.rnd.Read(b); err != nil {
		panic(err) // should never happen
	}
}

// randnulterm generates a random nul-terminated string of length in [a,b]
// inclusive of a and b, not including the nul byte.
func (g *generator) randnulterm(a, b int) []byte {
	const alpha = "abcdefghijklmnopqrstuvwxyz"
	n := g.randbtw(a, b)
	s := make([]byte, n+1)
	for i := 0; i < n; i++ {
		s[i] = alpha[g.rnd.Intn(len(alpha))]
	}
	return s
}

func intbytes(x int64, bits uint) []byte {
	n := bits / 8
	b := make([]byte, n)
	for i := uint(0); i < n; i++ {
		b[n-1-i] = byte(x)
		x >>= 8
	}
	return b
}
