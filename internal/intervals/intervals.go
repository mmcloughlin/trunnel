// Package intervals provides tools for manipulating collections of integer intervals.
package intervals

import (
	"fmt"
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

// Size returns the interval size.
func (i Interval) Size() uint64 {
	return i.hi - i.lo + 1
}

// Single returns true if the interval contains one integer.
func (i Interval) Single() bool {
	return i.Size() == 1
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
type Set []Interval

func (s Set) String() string {
	is := []string{}
	for _, i := range s {
		is = append(is, i.String())
	}
	return strings.Join(is, ",")
}

// Overlaps returns true if any intervals overlap.
func (s Set) Overlaps() bool {
	es := []edge{}
	for _, i := range s {
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
