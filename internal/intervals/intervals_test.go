package intervals

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/mmcloughlin/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRange(t *testing.T) {
	assert.Equal(t, Interval{lo: 10, hi: 20}, Range(10, 20))
}

func TestBadRange(t *testing.T) {
	assert.Panics(t, func() { Range(4, 3) })
}

func TestSingle(t *testing.T) {
	assert.Equal(t, Interval{lo: 13, hi: 13}, Single(13))
}

func TestBits(t *testing.T) {
	hi := ^uint64(0)
	for n := uint(0); n <= 64; n++ {
		assert.Equal(t, Interval{lo: 0, hi: hi >> n}, Bits(64-n))
	}
}

func TestIntervalSize(t *testing.T) {
	assert.Equal(t, uint64(1), Single(13).Size())
	assert.Equal(t, uint64(5), Range(13, 17).Size())
}

func TestIntervalString(t *testing.T) {
	assert.Equal(t, "13", Single(13).String())
	assert.Equal(t, "13-17", Range(13, 17).String())
}

func TestIntervalContains(t *testing.T) {
	cases := []struct {
		Interval Interval
		X        uint64
		Expect   bool
	}{
		{Single(42), 42, true},
		{Single(42), 100, false},
		{Range(100, 200), 150, true},
		{Range(100, 200), 100, true},
		{Range(100, 200), 200, true},
		{Range(100, 200), 99, false},
		{Range(100, 200), 201, false},
	}
	for _, c := range cases {
		assert.Equal(t, c.Expect, c.Interval.Contains(c.X))
	}
}

func TestIntType(t *testing.T) {
	assert.Equal(t, NewSet(Range(0, 127)), IntType(7))
}

func TestOverlaps(t *testing.T) {
	cases := []struct {
		Intervals []Interval
		Expect    bool
	}{
		{[]Interval{}, false},
		{[]Interval{Single(1)}, false},
		{[]Interval{Range(10, 100)}, false},
		{[]Interval{Range(10, 100), Range(50, 60)}, true},
		{[]Interval{Range(5, 10), Range(10, 15)}, true},
		{[]Interval{Range(5, 10), Range(11, 15)}, false},
	}
	for _, c := range cases {
		s := NewSet(c.Intervals...)
		t.Run(s.String(), func(t *testing.T) {
			assert.Equal(t, c.Expect, Overlaps(c.Intervals))
		})
	}
}

func TestSetContains(t *testing.T) {
	cases := []struct {
		Intervals []Interval
		X         uint64
		Expect    bool
	}{
		{[]Interval{}, 10, false},
		{[]Interval{Single(1)}, 1, true},
		{[]Interval{Single(1), Range(4, 5)}, 4, true},
	}
	for _, c := range cases {
		s := NewSet(c.Intervals...)
		t.Run(s.String(), func(t *testing.T) {
			assert.Equal(t, c.Expect, s.Contains(c.X))
		})
	}
}

func TestSetString(t *testing.T) {
	cases := []struct {
		Intervals []Interval
		Expect    string
	}{
		{
			Intervals: []Interval{},
			Expect:    "",
		},
		{
			Intervals: []Interval{Single(2), Range(4, 5)},
			Expect:    "2,4-5",
		},
		{
			Intervals: []Interval{Single(2), Range(4, 50), Range(30, 300)},
			Expect:    "2,4-300",
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.Expect, NewSet(c.Intervals...).String())
	}
}

func TestSubtract(t *testing.T) {
	cases := []struct {
		A      *Set
		B      *Set
		Expect *Set
	}{
		{
			NewSet(Range(1, 10)),
			NewSet(Range(5, 15)),
			NewSet(Range(1, 4)),
		},
		{
			NewSet(Range(50, 100)),
			NewSet(Range(25, 75)),
			NewSet(Range(76, 100)),
		},
		{
			NewSet(Range(0, 4), Range(8, 12)),
			NewSet(Range(2, 10)),
			NewSet(Range(0, 1), Range(11, 12)),
		},
		{
			NewSet(Range(50, 100)),
			NewSet(Single(75)),
			NewSet(Range(50, 74), Range(76, 100)),
		},
		{
			IntType(8),
			IntType(4),
			NewSet(Range(16, 255)),
		},
		{
			IntType(64),
			IntType(32),
			NewSet(Range(1<<32, math.MaxUint64)),
		},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("(%s)-(%s)", c.A, c.B), func(t *testing.T) {
			c.A.Subtract(c.B)
			assert.Equal(t, c.Expect, c.A)
		})
	}
}

func TestComplement(t *testing.T) {
	cases := []struct {
		Intervals []Interval
		Expect    []Interval
	}{
		{},
		{
			Intervals: []Interval{Range(10, 20)},
			Expect:    []Interval{OpenLeft(9), OpenRight(21)},
		},
		{
			Intervals: []Interval{OpenLeft(42)},
			Expect:    []Interval{OpenRight(43)},
		},
		{
			Intervals: []Interval{OpenRight(42)},
			Expect:    []Interval{OpenLeft(41)},
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.Expect, complement(c.Intervals))
	}
}

func TestSetRandomContains(t *testing.T) {
	s := NewSet(Single(1), Range(42, 53), Range(100, 1000))
	for i := 0; i < NumTrials(); i++ {
		require.True(t, s.Contains(s.Random()))
	}
}

func TestSetRandomObserveAll(t *testing.T) {
	s := NewSet(Single(1), Range(4200, 4201), Range(7, 10))
	counts := map[uint64]int{}
	for i := 0; i < NumTrials(); i++ {
		r := s.Random()
		if _, ok := counts[r]; !ok {
			counts[r] = 0
		}
		counts[r]++
	}
	t.Log(counts)

	expect := []uint64{1, 7, 8, 9, 10, 4200, 4201}
	for _, e := range expect {
		assert.Contains(t, counts, e)
	}
}

func TestSetRandomEmpty(t *testing.T) {
	assert.Panics(t, func() { Set{}.Random() })
}

func TestRandUint64n(t *testing.T) {
	rnd, err := random.NewFromSeeder(random.CryptoSeeder)
	require.NoError(t, err)
	for i := 0; i < NumTrials(); i++ {
		n := uint64(2 + rand.Intn(42))
		require.True(t, randuint64n(rnd, n) < n)
	}
}

func NumTrials() int {
	if testing.Short() {
		return 1000
	}
	return 100000
}
