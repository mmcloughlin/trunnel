package tv

import (
	"errors"

	"github.com/mmcloughlin/random"
	"github.com/mmcloughlin/trunnel/ast"
	"github.com/mmcloughlin/trunnel/fault"
	"github.com/mmcloughlin/trunnel/inspect"
	"github.com/mmcloughlin/trunnel/internal/intervals"
)

// Vector is a test vector.
type Vector struct {
	Data        []byte
	Constraints map[string]int64
}

// NewEmptyVector builds an empty test vector.
func NewEmptyVector() Vector {
	return Vector{
		Data:        []byte{},
		Constraints: map[string]int64{},
	}
}

type generator struct {
	structs     map[string]*ast.Struct
	resolver    *inspect.Resolver
	constraints map[string]int64
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
		v[s.Name], err = g.structure(s)
		if err != nil {
			return nil, err
		}
	}

	return v, nil
}

func (g *generator) structure(s *ast.Struct) ([]Vector, error) {
	n := len(s.Members)
	vectors := []Vector{NewEmptyVector()}
	for i := n - 1; i >= 0; i-- {
		extended := []Vector{}
		for _, v := range vectors {
			g.constraints = v.Constraints
			mvs, err := g.member(s.Members[i])
			if err != nil {
				return nil, err
			}

			for _, mv := range mvs {
				mv.Data = append(mv.Data, v.Data...)
				extended = append(extended, mv)
			}
		}
		vectors = extended
	}
	return vectors, nil
}

func (g *generator) member(m ast.Member) ([]Vector, error) {
	switch m := m.(type) {
	case *ast.Field:
		return g.field(m)
	default:
		return nil, fault.NewUnexpectedType(m)
	}
}

func (g *generator) field(f *ast.Field) ([]Vector, error) {
	switch t := f.Type.(type) {
	case *ast.IntType:
		return g.intType(f.Name, t)

	case *ast.NulTermString:
		return []Vector{
			{
				Data:        g.randnulterm(2, 20),
				Constraints: g.constraints,
			},
		}, nil

	case *ast.StructRef:
		s, ok := g.structs[t.Name]
		if !ok {
			return nil, errors.New("could not resolve struct name")
		}
		return g.structure(s)

	default:
		return nil, fault.NewUnexpectedType(t)
	}
}

func (g *generator) intType(name string, t *ast.IntType) ([]Vector, error) {
	var b []byte

	switch x, ok := g.constraints[name]; {
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

	return []Vector{
		{
			Data:        b,
			Constraints: g.constraints,
		},
	}, nil
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
	g.rnd.Read(b)
	return b
}

// randnulterm generates a random nul-terminated string of length in [a,b]
// inclusive of a and b, not including the nul byte.
func (g *generator) randnulterm(a, b int) []byte {
	const alpha = "abcdefghijklmnopqrstuvwxyz"
	n := a + g.rnd.Intn(b-a+1)
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
