// Package intervals provides tools for manipulating collections of integer intervals.
package intervals

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

// Interval represents the inclusive range of integers [lo, hi].
type Interval struct {
	lo uint64
	hi uint64
}

// Range builds the interval [l, h].
func Range(l, h uint64) Interval {
	if h < l {
		panic("bad range")
	}
	return Interval{lo: l, hi: h}
}

// Single builds the interval containing only x.
func Single(x uint64) Interval {
	return Range(x, x)
}

// Bits returns the interval [0, 2^n-1].
func Bits(n uint) Interval {
	return Range(0, (1<<n)-1)
}

// Size returns the interval size.
func (i Interval) Size() uint64 {
	return i.hi - i.lo + 1
}

// Single returns true if the interval contains one integer.
func (i Interval) Single() bool {
	return i.Size() == 1
}

// Contains returns whether x is contained in the interval.
func (i Interval) Contains(x uint64) bool {
	return i.lo <= x && x <= i.hi
}

func (i Interval) String() string {
	switch {
	case i.Single():
		return strconv.FormatUint(i.lo, 10)
	default:
		return fmt.Sprintf("%d-%d", i.lo, i.hi)
	}
}

// Set is a collection of intervals.
type Set struct {
	intervals []Interval
}

// NewSet builds a set from the given intervals.
func NewSet(is ...Interval) *Set {
	return &Set{intervals: is}
}

// IntType returns the set of possible values of an n-bit integer.
func IntType(n uint) *Set {
	return NewSet(Bits(n))
}

func (s Set) String() string {
	is := []string{}
	for _, i := range s.intervals {
		is = append(is, i.String())
	}
	return strings.Join(is, ",")
}

// Contains returns whether x is contained in the set.
func (s Set) Contains(x uint64) bool {
	for _, i := range s.intervals {
		if i.Contains(x) {
			return true
		}
	}
	return false
}

// Overlaps returns true if any intervals overlap.
func (s Set) Overlaps() bool {
	es := []edge{}
	for _, i := range s.intervals {
		es = append(es, edge{x: i.lo, d: 1})
		es = append(es, edge{x: i.hi, d: -1})
	}
	sort.Sort(edges(es))
	inside := 0
	for _, e := range es {
		inside += e.d
		if inside > 1 {
			return true
		}
	}
	return false
}

// Subtract subtracts other from s.
func (s *Set) Subtract(other *Set) {
	es := []edge{}
	for _, i := range s.intervals {
		es = append(es, edge{x: i.lo, d: 1})
		es = append(es, edge{x: i.hi, d: -1})
	}
	for _, i := range other.intervals {
		es = append(es, edge{x: i.lo - 1, d: -1})
		es = append(es, edge{x: i.hi + 1, d: 1})
	}
	sort.Sort(edges(es))
	n := 1
	inside := false
	result := []Interval{}
	var start uint64
	for _, e := range es {
		n += e.d
		if !inside && n >= 2 {
			start = e.x
			inside = true
		} else if inside && n < 2 {
			result = append(result, Range(start, e.x))
			inside = false
		}
	}
	s.intervals = result
}

type edge struct {
	x uint64
	d int
}

type edges []edge

func (e edges) Len() int      { return len(e) }
func (e edges) Swap(i, j int) { e[i], e[j] = e[j], e[i] }

func (e edges) Less(i, j int) bool {
	a, b := e[i], e[j]
	return a.x < b.x || (a.x == b.x && b.d < a.d)
}

// Random returns a random element of the collection. Assumes the collection
// contains non-overlapping intervals. Panics if s is empty.
func (s Set) Random() uint64 {
	if len(s.intervals) == 0 {
		panic("empty set")
	}
	type step struct {
		upper uint64
		delta uint64
	}
	steps := []step{}
	var cuml uint64
	for _, i := range s.intervals {
		cuml += i.Size()
		steps = append(steps, step{
			upper: cuml,
			delta: i.hi - cuml + 1,
		})
	}
	r := randuint64n(cuml)
	for _, step := range steps {
		if r < step.upper {
			return r + step.delta
		}
	}
	panic("unreachable")
}

// randuint64n returns a random uint64 in [0,n).
func randuint64n(n uint64) uint64 {
	mask := ^uint64(0)
	for mask > n {
		mask >>= 1
	}
	mask = (mask << 1) | uint64(1)

	for {
		r := randuint64() & mask
		if r < n {
			return r
		}
	}
}

func randuint64() uint64 {
	return uint64(rand.Int63())>>31 | uint64(rand.Int63())<<32
}
