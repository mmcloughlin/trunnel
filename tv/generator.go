package tv

import (
	"github.com/mmcloughlin/random"
	"github.com/mmcloughlin/trunnel/ast"
	"github.com/mmcloughlin/trunnel/fault"
	"github.com/mmcloughlin/trunnel/inspect"
	"github.com/mmcloughlin/trunnel/internal/intervals"
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
	structs     map[string]*ast.Struct
	resolver    *inspect.Resolver
	constraints Constraints
	rnd         random.Interface
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

func (g *generator) init(f *ast.File) error {
	s, err := inspect.Structs(f)
	if err != nil {
		return err
	}
	g.structs = s

	c, err := inspect.Constants(f)
	if err != nil {
		return err
	}
	g.resolver = inspect.NewResolver(c)

	return nil
}

func (g *generator) file(f *ast.File) (map[string][]Vector, error) {
	err := g.init(f)
	if err != nil {
		return nil, err
	}

	v := map[string][]Vector{}
	for _, s := range f.Structs {
		g.constraints = NewConstraints()
		v[s.Name], err = g.structure(s)
		if err != nil {
			return nil, err
		}
	}

	return v, nil
}

func (g *generator) structure(s *ast.Struct) ([]Vector, error) {
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

	case *ast.NulTermString:
		return []Vector{g.vector(g.randnulterm(2, 20))}, nil

	case *ast.StructRef:
		s, ok := g.structs[t.Name]
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
		s, err := g.intervals(t.Constraint)
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
	iv := g.vector([]byte{})
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
	// TODO(mbm): test vectors for length-constrained unions
	if u.Length != nil {
		return nil, fault.ErrNotImplemented
	}

	base := g.constraints.Clone()
	results := []Vector{}

	for _, c := range u.Cases {
		// TODO(mbm): test vectors union with default case
		if c.Case == nil {
			return nil, fault.ErrNotImplemented
		}

		i, err := g.intervals(c.Case)
		if err != nil {
			return nil, err
		}

		t := i.Random()
		g.constraints = base.Clone()
		g.constraints.SetRef(u.Tag, int64(t)) // XXX cast
		vs, err := g.members([]Vector{g.vector([]byte{})}, c.Members)
		if err != nil {
			return nil, err
		}
		results = append(results, vs...)
	}

	return results, nil
}

// vector builds vector with the current constraints.
func (g *generator) vector(b []byte) Vector {
	return Vector{
		Data:        b,
		Constraints: g.constraints,
	}
}

// intervals builds intervals object from an integer list.
func (g *generator) intervals(l *ast.IntegerList) (intervals.Set, error) {
	s := make(intervals.Set, len(l.Ranges))
	for i, r := range l.Ranges {
		lo, err := g.resolver.Integer(r.Low)
		if err != nil {
			return nil, err
		}
		if r.High == nil {
			s[i] = intervals.Single(uint64(lo)) // XXX cast
			continue
		}
		hi, err := g.resolver.Integer(r.High)
		if err != nil {
			return nil, err
		}
		s[i] = intervals.Range(uint64(lo), uint64(hi)) // XXX cast
	}
	if s.Overlaps() {
		return nil, errors.New("overlapping intervals")
	}
	return s, nil
}

func (g *generator) randint(bits int) []byte {
	n := bits / 8
	b := make([]byte, n)
	if _, err := g.rnd.Read(b); err != nil {
		panic(err) // should never happen
	}
	return b
}

// randbtw generates an integer between a and b, inclusive.
func (g *generator) randbtw(a, b int) int {
	return a + g.rnd.Intn(b-a+1)
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

func intbytes(x int64, bits int) []byte {
	n := bits / 8
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[n-1-i] = byte(x)
		x >>= 8
	}
	return b
}
